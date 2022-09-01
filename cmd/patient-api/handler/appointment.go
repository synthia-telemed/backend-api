package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/server/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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
	logger           *zap.SugaredLogger
	*server.GinHandler
}

func NewAppointmentHandler(patientDS datastore.PatientDataStore, paymentDS datastore.PaymentDataStore, hos hospital.SystemClient, logger *zap.SugaredLogger) *AppointmentHandler {
	return &AppointmentHandler{
		patientDataStore: patientDS,
		hospitalClient:   hos,
		paymentDataStore: paymentDS,
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
	apps, err := h.hospitalClient.ListAppointmentsByPatientID(context.Background(), patient.RefID)
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
	c.JSON(http.StatusOK, res)
}

type GetAppointmentResponse struct {
	*hospital.Appointment
	Payment *datastore.Payment `json:"payment"`
}

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
