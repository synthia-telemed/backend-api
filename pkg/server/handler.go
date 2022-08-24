package server

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler interface {
	Register(r *gin.RouterGroup)
}

type GinHandler struct {
	Logger *zap.SugaredLogger
}

func (h GinHandler) InternalServerError(c *gin.Context, err error, msg string) {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.CaptureException(err)
	}
	h.Logger.Errorw(msg, "error", err.Error())
	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
}

func (h GinHandler) GetUserID(c *gin.Context) uint {
	id, _ := c.Get("UserID")
	return id.(uint)
}
