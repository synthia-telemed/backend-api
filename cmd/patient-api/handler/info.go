package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
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
	g := r.Group("info", h.ParseUserID, h.ParsePatient, h.ParseHospitalPatientInfo)
	g.GET("", h.GetPatientInfo)
	g.GET("/name", h.GetName)
}

type GetNameResponse struct {
	EN *hospital.Name `json:"EN"`
	TH *hospital.Name `json:"TH"`
}

// GetName godoc
// @Summary      Get patient name
// @Tags         Info
// @Success      200  {object}	GetNameResponse "Name of the patient in both Thai and English"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      404  {object}  server.ErrorResponse "Patient not found"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /info/name [get]
func (h InfoHandler) GetName(c *gin.Context) {
	rawInfo, _ := c.Get("PatientInfo")
	patientInfo := rawInfo.(*hospital.Patient)
	c.JSON(http.StatusOK, &GetNameResponse{
		EN: patientInfo.NameEN,
		TH: patientInfo.NameTH,
	})
}

// GetPatientInfo godoc
// @Summary      Get patient information
// @Tags         Info
// @Success      200  {object}	hospital.Patient "Patient information"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      404  {object}  server.ErrorResponse "Patient not found"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /info [get]
func (h InfoHandler) GetPatientInfo(c *gin.Context) {
	rawInfo, _ := c.Get("PatientInfo")
	patientInfo := rawInfo.(*hospital.Patient)
	c.JSON(http.StatusOK, patientInfo)
}

func (h InfoHandler) ParseHospitalPatientInfo(c *gin.Context) {
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
		c.AbortWithStatusJSON(http.StatusNotFound, ErrPatientNotFound)
		return
	}
	c.Set("PatientInfo", patientInfo)
}
