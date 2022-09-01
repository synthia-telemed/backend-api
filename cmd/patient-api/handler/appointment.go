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
)

type AppointmentHandler struct {
	patientDataStore datastore.PatientDataStore
	hospitalClient   hospital.SystemClient
	logger           *zap.SugaredLogger
	*server.GinHandler
}

func NewAppointmentHandler(patientDS datastore.PatientDataStore, hos hospital.SystemClient, logger *zap.SugaredLogger) *AppointmentHandler {
	return &AppointmentHandler{
		patientDataStore: patientDS,
		hospitalClient:   hos,
		logger:           logger,
		GinHandler:       &server.GinHandler{Logger: logger},
	}
}

func (h AppointmentHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/appointment")
	g.GET("", middleware.ParseUserID, h.ParsePatient, h.ListAppointments)
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
