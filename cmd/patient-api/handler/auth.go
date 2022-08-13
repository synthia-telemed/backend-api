package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	. "github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/sms"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AuthHandler struct {
	patientDataStore  PatientDataStore
	hospitalSysClient hospital.SystemClient
	smsClient         sms.Client
	cacheClient       cache.Client
	logger            *zap.SugaredLogger
}

func NewAuthHandler(patientDataStore PatientDataStore, hosClient hospital.SystemClient, sms sms.Client, cache cache.Client, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		patientDataStore:  patientDataStore,
		hospitalSysClient: hosClient,
		smsClient:         sms,
		cacheClient:       cache,
		logger:            logger,
	}
}

type LoginRequest struct {
	Credential string `json:"credential" binding:"required"`
}

type LoginResponse struct {
	PhoneNumber string `json:"phone_number"`
}

func (h AuthHandler) Signin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	patientInfo, err := h.hospitalSysClient.FindPatientByGovCredential(context.Background(), req.Credential)
	if err != nil {
		InternalServerError(c, h.logger, err, "h.hospitalSysClient.FindPatientByGovCredential error")
		return
	}
	if patientInfo == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{"Patient not found"})
		return
	}

	otp, err := gonanoid.Generate("1234567890", 6)
	if err != nil {
		InternalServerError(c, h.logger, err, "gonanoid.Generate error")
		return
	}

	if err := h.cacheClient.Set(context.Background(), otp, patientInfo.Id, time.Minute*10); err != nil {
		InternalServerError(c, h.logger, err, "h.cacheClient.Set error")
		return
	}
	if err := h.smsClient.Send(patientInfo.PhoneNumber, fmt.Sprintf("Your OTP is %s", otp)); err != nil {
		InternalServerError(c, h.logger, err, "h.smsClient.Send error")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"phone_number": h.censorPhoneNumber(patientInfo.PhoneNumber),
	})
}

func (h AuthHandler) censorPhoneNumber(number string) string {
	return number[:3] + "***" + number[len(number)-4:]
}
