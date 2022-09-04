package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/clock"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/server/middleware"
	"go.uber.org/zap"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var (
	ErrAppointmentIDMissing = server.NewErrorResponse("Appointment ID is missing")
	ErrAppointmentIDInvalid = server.NewErrorResponse("Invalid appointment ID")
	ErrAppointmentNotFound  = server.NewErrorResponse("Appointment not found")
	ErrForbidden            = server.NewErrorResponse("Forbidden")
)

type AppointmentHandler struct {
	patientDataStore datastore.PatientDataStore
	paymentDataStore datastore.PaymentDataStore
	hospitalClient   hospital.SystemClient
	clock            clock.Clock
	logger           *zap.SugaredLogger
	*server.GinHandler
}

func NewAppointmentHandler(patientDS datastore.PatientDataStore, paymentDS datastore.PaymentDataStore, hos hospital.SystemClient, c clock.Clock, logger *zap.SugaredLogger) *AppointmentHandler {
	return &AppointmentHandler{
		patientDataStore: patientDS,
		hospitalClient:   hos,
		paymentDataStore: paymentDS,
		clock:            c,
		logger:           logger,
		GinHandler:       &server.GinHandler{Logger: logger},
	}
}

func (h AppointmentHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/appointment")
	g.GET("", middleware.ParseUserID, h.ParsePatient, h.ListAppointments)
	g.GET("/:appointmentID", middleware.ParseUserID, h.ParsePatient, h.GetAppointment)
}

type ListAppointmentsResponse struct {
	Completed []*hospital.AppointmentOverview `json:"completed"`
	Scheduled []*hospital.AppointmentOverview `json:"scheduled"`
	Cancelled []*hospital.AppointmentOverview `json:"cancelled"`
}

// ListAppointments godoc
// @Summary      Get list of appointment of the patient
// @Tags         Appointment
// @Success      200  {object}	ListAppointmentsResponse "List of appointment group by status"
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
	res := ListAppointmentsResponse{}
	for _, a := range apps {
		switch a.Status {
		case hospital.AppointmentStatusCancelled:
			res.Cancelled = append(res.Cancelled, a)
		case hospital.AppointmentStatusCompleted:
			res.Completed = append(res.Completed, a)
		case hospital.AppointmentStatusScheduled:
			res.Scheduled = append(res.Scheduled, a)
		}
	}
	ReverseSlice(res.Scheduled)
	c.JSON(http.StatusOK, res)
}

type GetAppointmentResponse struct {
	*hospital.Appointment
	Payment *datastore.Payment `json:"payment"`
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
	apps, err := h.hospitalClient.FindAppointmentByID(context.Background(), int(appointmentID))
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalClient.FindAppointmentByID error")
		return
	}
	if apps == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrAppointmentNotFound)
		return
	}
	if apps.PatientID != patient.RefID {
		c.AbortWithStatusJSON(http.StatusForbidden, ErrForbidden)
	}
	res := &GetAppointmentResponse{
		Appointment: apps,
		Payment:     nil,
	}
	if apps.Status == hospital.AppointmentStatusCompleted {
		payment, err := h.paymentDataStore.FindLatestByInvoiceIDAndStatus(apps.Invoice.Id, datastore.SuccessPaymentStatus)
		if err != nil {
			h.InternalServerError(c, err, "h.paymentDataStore.FindByInvoiceID error")
			return
		}
		res.Payment = payment
	}
	c.JSON(http.StatusOK, res)
}

func (h AppointmentHandler) ParsePatient(c *gin.Context) {
	patientID := h.GetUserID(c)
	patient, err := h.patientDataStore.FindByID(patientID)
	if err != nil {
		h.InternalServerError(c, err, "h.patientDataStore.FindByID error")
		return
	}
	if patient == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrPatientNotFound)
		return
	}
	c.Set("Patient", patient)
}

func ReverseSlice[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}
