package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/clock"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/sms"
	"github.com/synthia-telemed/backend-api/pkg/token"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	ErrInvalidRequestBody = server.NewErrorResponse("Invalid request body")
	ErrPatientNotFound    = server.NewErrorResponse("Patient not found")
	ErrInvalidOTP         = server.NewErrorResponse("OTP is invalid or expired")
)

type AuthHandler struct {
	patientDataStore  datastore.PatientDataStore
	hospitalSysClient hospital.SystemClient
	smsClient         sms.Client
	cacheClient       cache.Client
	tokenService      token.Service
	clock             clock.Clock
	PatientGinHandler
}

func NewAuthHandler(patientDataStore datastore.PatientDataStore, hosClient hospital.SystemClient, sms sms.Client, cache cache.Client, tokenService token.Service, clock clock.Clock, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		patientDataStore:  patientDataStore,
		hospitalSysClient: hosClient,
		smsClient:         sms,
		cacheClient:       cache,
		tokenService:      tokenService,
		clock:             clock,
		PatientGinHandler: NewPatientGinHandler(patientDataStore, logger),
	}
}

func (h AuthHandler) Register(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	authGroup.POST("/signin", h.Signin)
	authGroup.POST("/verify", h.VerifyOTP)
	authGroup.DELETE("/signout", h.ParseUserID, h.ParsePatient, h.SignOut)
}

type SigninRequest struct {
	Credential string `json:"credential" binding:"required"`
}

type SigninResponse struct {
	ExpiredAt   time.Time `json:"expired_at"`
	PhoneNumber string    `json:"phone_number"`
}

// Signin godoc
// @Summary      Start signing-in with government credential
// @Description  Initiate auth process with government credential which will sent OTP to patient's phone number
// @Tags         Auth
// @Param 	  	 SigninRequest body SigninRequest true "Patient government credential (Passport ID or National ID)"
// @Success      201  {object}  SigninResponse "OTP is sent to patient's phone number"
// @Failure      400  {object}  server.ErrorResponse "Invalid request body"
// @Failure      404  {object}  server.ErrorResponse "Provided credential is not in the hospital system"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Router       /auth/signin [post]
func (h AuthHandler) Signin(c *gin.Context) {
	var req SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	patientInfo, err := h.hospitalSysClient.FindPatientByGovCredential(context.Background(), req.Credential)
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalSysClient.FindPatientByGovCredential error")
		return
	}
	if patientInfo == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrPatientNotFound)
		return
	}

	otp, err := gonanoid.Generate("1234567890", 6)
	if err != nil {
		h.InternalServerError(c, err, "gonanoid.Generate error")
		return
	}

	expiredIn := time.Minute * 10
	if err := h.cacheClient.Set(context.Background(), otp, patientInfo.Id, expiredIn); err != nil {
		h.InternalServerError(c, err, "h.cacheClient.Set error")
		return
	}
	if err := h.smsClient.Send(patientInfo.PhoneNumber, fmt.Sprintf("Your OTP is %s", otp)); err != nil {
		h.InternalServerError(c, err, "h.smsClient.Send error")
		return
	}

	c.JSON(http.StatusCreated, &SigninResponse{
		PhoneNumber: h.censorPhoneNumber(patientInfo.PhoneNumber),
		ExpiredAt:   h.clock.Now().Add(expiredIn),
	})
}

type VerifyOTPRequest struct {
	OTP string `json:"otp" binding:"required"`
}

type VerifyOTPResponse struct {
	Token string `json:"token"`
}

// VerifyOTP godoc
// @Summary      Verify OTP and get token
// @Description  Complete auth process with OTP verification. It will return token if verification success.
// @Tags         Auth
// @Param 	  	 VerifyOTPRequest body VerifyOTPRequest true "OTP that is sent to patient's phone number"
// @Success      201  {object}  VerifyOTPResponse "JWS Token for later use"
// @Failure      400  {object}  server.ErrorResponse "Invalid request body"
// @Failure      400  {object}  server.ErrorResponse "OTP is invalid or expired"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Router       /auth/verify [post]
func (h AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	refID, err := h.cacheClient.Get(context.Background(), req.OTP, true)
	if err != nil {
		h.InternalServerError(c, err, "h.cacheClient.Get error")
		return
	}
	if len(refID) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidOTP)
		return
	}

	patient := &datastore.Patient{RefID: refID}
	if err := h.patientDataStore.FindOrCreate(patient); err != nil {
		h.InternalServerError(c, err, "h.patientDataStore.FindByRefID error")
		return
	}

	jws, err := h.tokenService.GenerateToken(uint64(patient.ID), "Patient")
	if err != nil {
		h.InternalServerError(c, err, "h.tokenService.GenerateToken error")
		return
	}
	c.JSON(http.StatusCreated, VerifyOTPResponse{Token: jws})
}

func (h AuthHandler) censorPhoneNumber(number string) string {
	return number[:3] + "***" + number[len(number)-4:]
}

func (h AuthHandler) SignOut(c *gin.Context) {
	patientRaw, _ := c.Get("Patient")
	patient := patientRaw.(*datastore.Patient)

	patient.NotificationToken = ""
	if err := h.patientDataStore.Save(patient); err != nil {
		h.InternalServerError(c, err, "h.patientDataStore.Save")
		return
	}
	c.AbortWithStatus(http.StatusOK)
}
