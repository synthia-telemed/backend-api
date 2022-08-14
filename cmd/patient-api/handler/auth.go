package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/sms"
	"github.com/synthia-telemed/backend-api/pkg/token"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AuthHandler struct {
	patientDataStore  datastore.PatientDataStore
	hospitalSysClient hospital.SystemClient
	smsClient         sms.Client
	cacheClient       cache.Client
	logger            *zap.SugaredLogger
	tokenService      token.Service
}

func NewAuthHandler(patientDataStore datastore.PatientDataStore, hosClient hospital.SystemClient, sms sms.Client, cache cache.Client, tokenService token.Service, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		patientDataStore:  patientDataStore,
		hospitalSysClient: hosClient,
		smsClient:         sms,
		cacheClient:       cache,
		logger:            logger,
		tokenService:      tokenService,
	}
}

type SigninRequest struct {
	Credential string `json:"credential" binding:"required"`
}

type SigninResponse struct {
	PhoneNumber string `json:"phone_number"`
}

func (h AuthHandler) Signin(c *gin.Context) {
	var req SigninRequest
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

type VerifyOTPRequest struct {
	OTP string `json:"otp" binding:"required"`
}

type VerifyOTPResponse struct {
	Token string `json:"token"`
}

func (h AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	refID, err := h.cacheClient.Get(context.Background(), req.OTP)
	if err != nil {
		InternalServerError(c, h.logger, err, "h.cacheClient.Get error")
		return
	}
	if len(refID) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{"OTP is invalid or expired"})
		return
	}

	patient, err := h.patientDataStore.FindByRefID(refID)
	if err != nil {
		InternalServerError(c, h.logger, err, "h.patientDataStore.FindByRefID error")
		return
	}
	if patient == nil {
		patient = &datastore.Patient{RefID: refID}
		if err := h.patientDataStore.Create(patient); err != nil {
			InternalServerError(c, h.logger, err, "h.patientDataStore.Create error")
			return
		}
	}
	h.logger.Info(patient)

	jws, err := h.tokenService.GenerateToken(uint64(patient.ID), "Patient")
	if err != nil {
		InternalServerError(c, h.logger, err, "h.tokenService.GenerateToken error")
		return
	}

	c.JSON(http.StatusCreated, VerifyOTPResponse{Token: jws})
}

func (h AuthHandler) censorPhoneNumber(number string) string {
	return number[:3] + "***" + number[len(number)-4:]
}
