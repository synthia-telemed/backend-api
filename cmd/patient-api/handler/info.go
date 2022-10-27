package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/server/middleware"
	"go.uber.org/zap"
	"net/http"
)

type InfoHandler struct {
	hospitalClient hospital.SystemClient
	PatientGinHandler
}

func NewInfoHandler(patientDataStore datastore.PatientDataStore, hospitalClient hospital.SystemClient, logger *zap.SugaredLogger) *InfoHandler {
	return &InfoHandler{
		hospitalClient:    hospitalClient,
		PatientGinHandler: NewPatientGinHandler(patientDataStore, logger),
	}
}

func (h InfoHandler) Register(r *gin.RouterGroup) {
	g := r.Group("info")
	g.GET("/name", middleware.ParseUserID, h.ParsePatient, h.GetName)
}

type GetNameResponse struct {
	EN *hospital.Name `json:"EN"`
	TH *hospital.Name `json:"TH"`
}

func (h InfoHandler) GetName(c *gin.Context) {
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
	patientInfo, err := h.hospitalClient.FindPatientByID(context.Background(), patient.RefID)
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalClient.FindPatientByID error")
		return
	}
	if patientInfo == nil {
		c.JSON(http.StatusNotFound, ErrPatientNotFound)
		return
	}
	c.JSON(http.StatusOK, &GetNameResponse{
		EN: patientInfo.NameEN,
		TH: patientInfo.NameTH,
	})
}
