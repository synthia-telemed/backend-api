package handler

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	. "github.com/synthia-telemed/backend-api/pkg/datastore"
	"net/http"
)

type AuthHandler struct {
	patientDataStore PatientDataStore
}

func NewAuthHandler(patientDataStore PatientDataStore) *AuthHandler {
	return &AuthHandler{patientDataStore}
}

func (h AuthHandler) Signin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	patient, err := h.patientDataStore.FindByNationalID(req.NationalID)
	if err != nil {
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureException(err)
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	if patient == nil {
		// TODO: Query patient from hospital system
		// TODO: Create new patient
	}

	// TODO: send OTP to patient's phone number
	c.Status(201)
}

type LoginRequest struct {
	NationalID string `json:"national_id" binding:"required,len=13"`
}

type LoginResponse struct {
	PhoneNumber string `json:"phone_number"`
}
