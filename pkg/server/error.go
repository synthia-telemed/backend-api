package server

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"time"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Message: msg}
}

func AssertFatalError(logger *zap.SugaredLogger, err error, msg string) {
	if err == nil {
		return
	}
	sentry.CaptureException(err)
	sentry.Flush(time.Second * 2)
	logger.Fatalw(msg, "error", err)
}
