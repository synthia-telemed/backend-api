package server

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var (
	ErrInvalidRequestBody = ErrorResponse{Message: "Invalid request body"}
	ErrPatientNotFound    = ErrorResponse{Message: "Patient not found"}
	ErrInvalidOTP         = ErrorResponse{Message: "OTP is invalid or expired"}
	// Payment

	ErrFailedToAddCreditCard = ErrorResponse{Message: "Failed to add credit card"}
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func InternalServerError(c *gin.Context, logger *zap.SugaredLogger, err error, message string) {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.CaptureException(err)
	}
	logger.Errorw(message, "error", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
}
