package handler

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	. "github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	patientDataStore  PatientDataStore
	hospitalSysClient hospital.SystemClient
	logger            *zap.SugaredLogger
}

func NewAuthHandler(patientDataStore PatientDataStore, hosClient hospital.SystemClient, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{patientDataStore, hosClient, logger}
}

func (h AuthHandler) Signin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	patient, err := h.patientDataStore.FindByGovCredential(req.Credential)
	if err != nil {
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureException(err)
		}
		h.logger.Errorw("h.patientDataStore.FindByGovCredential error", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	if patient == nil {
		patient, err := h.hospitalSysClient.FindPatientByGovCredential(req.Credential)
		if err != nil {
			if hub := sentrygin.GetHubFromContext(c); hub != nil {
				hub.CaptureException(err)
			}
			h.logger.Errorw("h.hospitalSysClient.FindPatientByGovCredential error", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}

		// TODO: Create new patient
		if err := h.patientDataStore.Create(patient); err != nil {
			if hub := sentrygin.GetHubFromContext(c); hub != nil {
				hub.CaptureException(err)
			}
			h.logger.Errorw("h.patientDataStore.Create error", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}
	}

	// TODO: send OTP to patient's phone number
	c.Status(201)
}

type LoginRequest struct {
	Credential string `json:"credential" binding:"required"`
}

type LoginResponse struct {
	PhoneNumber string `json:"phone_number"`
}
