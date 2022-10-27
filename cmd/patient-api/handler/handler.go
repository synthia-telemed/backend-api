package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"go.uber.org/zap"
	"net/http"
)

type PatientGinHandler struct {
	patientDataStore datastore.PatientDataStore
	server.GinHandler
}

func NewPatientGinHandler(patientDS datastore.PatientDataStore, logger *zap.SugaredLogger) PatientGinHandler {
	return PatientGinHandler{
		patientDataStore: patientDS,
		GinHandler:       server.GinHandler{Logger: logger},
	}
}

func (h PatientGinHandler) ParsePatient(c *gin.Context) {
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
