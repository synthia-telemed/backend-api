package server

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Message: msg}
}

func InternalServerError(c *gin.Context, logger *zap.SugaredLogger, err error, message string) {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.CaptureException(err)
	}
	logger.Errorw(message, "error", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
}
