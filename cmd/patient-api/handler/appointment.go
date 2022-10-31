package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/clock"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrAppointmentIDMissing = server.NewErrorResponse("Appointment ID is missing")
	ErrAppointmentIDInvalid = server.NewErrorResponse("Invalid appointment ID")
	ErrAppointmentNotFound  = server.NewErrorResponse("Appointment not found")
	ErrForbidden            = server.NewErrorResponse("Forbidden")
	ErrRoomIDNotFound       = server.NewErrorResponse("RoomID of the appointment not found")
)

type AppointmentHandler struct {
	patientDataStore     datastore.PatientDataStore
	paymentDataStore     datastore.PaymentDataStore
	appointmentDataStore datastore.AppointmentDataStore
	hospitalClient       hospital.SystemClient
	cacheClient          cache.Client
	clock                clock.Clock
	PatientGinHandler
}

func NewAppointmentHandler(patientDS datastore.PatientDataStore, paymentDS datastore.PaymentDataStore, appsDS datastore.AppointmentDataStore, hos hospital.SystemClient, cacheClient cache.Client, c clock.Clock, logger *zap.SugaredLogger) *AppointmentHandler {
	return &AppointmentHandler{
		patientDataStore:     patientDS,
		hospitalClient:       hos,
		paymentDataStore:     paymentDS,
		appointmentDataStore: appsDS,
		cacheClient:          cacheClient,
		clock:                c,
		PatientGinHandler:    NewPatientGinHandler(patientDS, logger),
	}
}

func (h AppointmentHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/appointment", h.ParseUserID, h.ParsePatient)
	g.GET("", h.ListAppointments)
	g.GET("/next", h.GetNextScheduledAppointment)
	g.GET("/:appointmentID", h.AuthorizedPatientToAppointment, h.GetAppointment)
	g.GET("/:appointmentID/roomID", h.AuthorizedPatientToAppointment, h.GetAppointmentRoomID)
}

// GetNextScheduledAppointment godoc
// @Summary      Get next scheduled appointment
// @Tags         Appointment
// @Success      200  {object}	hospital.AppointmentOverview "Overview of the appointment detail. If there is no scheduled appointment, empty body is returned"
// @Failure      400  {object}  server.ErrorResponse "Patient not found"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /appointment/next [get]
func (h AppointmentHandler) GetNextScheduledAppointment(c *gin.Context) {
	rawPatient, exist := c.Get("Patient")
	if !exist {
		h.InternalServerError(c, errors.New("c.Get Patient not exist"), "c.Get Patient not exist")
		return
	}
	patient, ok := rawPatient.(*datastore.Patient)
	if !ok {
		h.InternalServerError(c, errors.New("patient type casting error"), "Patient type casting error")
		return
	}
	appointments, err := h.hospitalClient.ListAppointmentsWithFilters(context.Background(), &hospital.ListAppointmentsFilters{
		PatientID: &patient.RefID,
		Status:    hospital.AppointmentStatusScheduled,
	}, 1, 0)
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalClient.ListAppointmentsWithFilters error")
		return
	}
	if len(appointments) == 0 {
		c.AbortWithStatus(http.StatusOK)
		return
	}
	c.JSON(http.StatusOK, appointments[0])
}

// ListAppointments godoc
// @Summary      Get list of appointment of the patient
// @Tags         Appointment
// @Success      200  {object}	hospital.CategorizedAppointment "List of appointment group by status"
// @Failure      400  {object}  server.ErrorResponse "Patient not found"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /appointment [get]
func (h AppointmentHandler) ListAppointments(c *gin.Context) {
	rawPatient, exist := c.Get("Patient")
	if !exist {
		h.InternalServerError(c, errors.New("c.Get Patient not exist"), "c.Get Patient not exist")
		return
	}
	patient, ok := rawPatient.(*datastore.Patient)
	if !ok {
		h.InternalServerError(c, errors.New("patient type casting error"), "Patient type casting error")
		return
	}
	since := h.clock.Now().Add(-time.Hour * 24 * 365 * 3) // 3 years
	apps, err := h.hospitalClient.ListAppointmentsByPatientID(context.Background(), patient.RefID, since)
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalClient.ListAppointmentsByPatientID error")
		return
	}
	c.JSON(http.StatusOK, h.hospitalClient.CategorizeAppointmentByStatus(apps))
}

type GetAppointmentResponse struct {
	*hospital.Appointment
	Payment  *datastore.Payment `json:"payment"`
	Duration *float64           `json:"duration"`
}

// GetAppointment godoc
// @Summary      Get an appointment detail by appointment ID
// @Tags         Appointment
// @Param  		 appointmentID 	path	 integer 	true "ID of the appointment"
// @Success      200  {object}	GetAppointmentResponse "An appointment detail"
// @Failure      400  {object}  server.ErrorResponse "Patient not found"
// @Failure      400  {object}  server.ErrorResponse "appointmentID is not provided"
// @Failure      400  {object}  server.ErrorResponse "appointmentID is invalid"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      403  {object}  server.ErrorResponse "The patient doesn't own the appointment"
// @Failure      404  {object}  server.ErrorResponse "Appointment not found"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /appointment/{appointmentID} [get]
func (h AppointmentHandler) GetAppointment(c *gin.Context) {
	rawAppointment, _ := c.Get("Appointment")
	appointment, _ := rawAppointment.(*hospital.Appointment)
	res := &GetAppointmentResponse{
		Appointment: appointment,
		Payment:     nil,
	}
	if appointment.Status == hospital.AppointmentStatusCompleted {
		compAppData, err := h.appointmentDataStore.FindByRefID(appointment.Id)
		if err != nil {
			h.InternalServerError(c, err, "h.appointmentDataStore.FindByRefID error")
			return
		}
		res.Duration = &compAppData.Duration
	}
	if appointment.Invoice != nil && appointment.Invoice.Paid {
		payment, err := h.paymentDataStore.FindLatestByInvoiceIDAndStatus(appointment.Invoice.Id, datastore.SuccessPaymentStatus)
		if err != nil {
			h.InternalServerError(c, err, "h.paymentDataStore.FindByInvoiceID error")
			return
		}
		res.Payment = payment
	}
	c.JSON(http.StatusOK, res)
}

type GetAppointmentRoomIDResponse struct {
	RoomID string `json:"room_id"`
}

// GetAppointmentRoomID godoc
// @Summary      Get room ID of the appointment
// @Tags         Appointment
// @Param  		 appointmentID 	path	 integer 	true "ID of the appointment"
// @Success      200  {object}	GetAppointmentRoomIDResponse "Room ID for the appointment"
// @Failure      400  {object}  server.ErrorResponse "Patient not found"
// @Failure      400  {object}  server.ErrorResponse "appointmentID is not provided"
// @Failure      400  {object}  server.ErrorResponse "appointmentID is invalid"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      403  {object}  server.ErrorResponse "The patient doesn't own the appointment"
// @Failure      404  {object}  server.ErrorResponse "Appointment not found"
// @Failure      404  {object}  server.ErrorResponse "RoomID of the appointment not found"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /appointment/{appointmentID}/roomID [get]
func (h AppointmentHandler) GetAppointmentRoomID(c *gin.Context) {
	rawAppointment, _ := c.Get("Appointment")
	appointment, _ := rawAppointment.(*hospital.Appointment)
	roomID, err := h.cacheClient.Get(context.Background(), cache.AppointmentRoomIDKey(appointment.Id), false)
	if err != nil {
		h.InternalServerError(c, err, "h.cacheClient.Get error")
		return
	}
	if roomID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrRoomIDNotFound)
		return
	}
	c.JSON(http.StatusOK, &GetAppointmentRoomIDResponse{RoomID: roomID})
}

func (h AppointmentHandler) AuthorizedPatientToAppointment(c *gin.Context) {
	rawPatient, exist := c.Get("Patient")
	if !exist {
		h.InternalServerError(c, errors.New("c.Get Patient not exist"), "c.Get Patient not exist")
		return
	}
	patient, ok := rawPatient.(*datastore.Patient)
	if !ok {
		h.InternalServerError(c, errors.New("patient type casting error"), "Patient type casting error")
		return
	}
	appointmentIDStr := c.Param("appointmentID")
	if appointmentIDStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrAppointmentIDMissing)
		return
	}
	appointmentID, err := strconv.ParseInt(appointmentIDStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrAppointmentIDInvalid)
		return
	}
	appointment, err := h.hospitalClient.FindAppointmentByID(context.Background(), int(appointmentID))
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalClient.FindAppointmentByID error")
		return
	}
	if appointment == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrAppointmentNotFound)
		return
	}
	if appointment.PatientID != patient.RefID {
		c.AbortWithStatusJSON(http.StatusForbidden, ErrForbidden)
		return
	}
	c.Set("Appointment", appointment)
}
