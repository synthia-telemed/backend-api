package handler

import (
	"context"
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid"
	. "github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/sms"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	patientDataStore  PatientDataStore
	hospitalSysClient hospital.SystemClient
	smsClient         sms.Client
	logger            *zap.SugaredLogger
}

func NewAuthHandler(patientDataStore PatientDataStore, hosClient hospital.SystemClient, sms sms.Client, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{patientDataStore, hosClient, sms, logger}
}

func (h AuthHandler) Signin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	patientInfo, err := h.hospitalSysClient.FindPatientByGovCredential(context.Background(), req.Credential)
	if err != nil {
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureException(err)
		}
		h.logger.Errorw("h.hospitalSysClient.FindPatientByGovCredential error", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	if patientInfo == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{"Patient not found"})
		return
	}

	// TODO: Create new patient
	//if err := h.patientDataStore.Create(&Patient{}); err != nil {
	//	if hub := sentrygin.GetHubFromContext(c); hub != nil {
	//		hub.CaptureException(err)
	//	}
	//	h.logger.Errorw("h.patientDataStore.Create error", "error", err)
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
	//	return
	//}

	// TODO: send OTP to patient's phone number
	otp, err := gonanoid.Generate("1234567890", 6)
	if err != nil {
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureException(err)
		}
		h.logger.Errorw("gonanoid.Generate error", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	if err := h.smsClient.Send(patientInfo.PhoneNumber, fmt.Sprintf("Your OTP is %s", otp)); err != nil {
		InternalServerError(c, h.logger, err, "h.smsClient.Send error")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"phone_number": patientInfo.PhoneNumber,
	})
}

type LoginRequest struct {
	Credential string `json:"credential" binding:"required"`
}

type LoginResponse struct {
	PhoneNumber string `json:"phone_number"`
}
