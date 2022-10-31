package server

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler interface {
	Register(r *gin.RouterGroup)
}

type GinHandler struct {
	Logger *zap.SugaredLogger
}

func (h GinHandler) InternalServerError(c *gin.Context, err error, msg string) {
	h.InternalServerErrorWithoutAborting(c, err, msg)
	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
}

func (h GinHandler) InternalServerErrorWithoutAborting(c *gin.Context, err error, msg string) {
	h.Logger.Errorw(msg, "error", err.Error())
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.CaptureException(err)
	}
}

func (h GinHandler) ParseUserID(c *gin.Context) {
	id := c.Request.Header.Get("X-USER-ID")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing user ID"})
		return
	}
	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid user ID"})
		return
	}
	c.Set("UserID", uint(uintID))
}

func (h GinHandler) GetUserID(c *gin.Context) uint {
	id, _ := c.Get("UserID")
	return id.(uint)
}
