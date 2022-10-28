package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"go.uber.org/zap"
	"net/http"
)

type NotificationHandler struct {
	notificationDataStore datastore.NotificationDataStore
	PatientGinHandler
}

func NewNotificationHandler(notificationDataStore datastore.NotificationDataStore, patientDataStore datastore.PatientDataStore, logger *zap.SugaredLogger) *NotificationHandler {
	return &NotificationHandler{
		notificationDataStore: notificationDataStore,
		PatientGinHandler:     NewPatientGinHandler(patientDataStore, logger),
	}
}

func (h NotificationHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/notification")
	g.GET("", h.ParseUserID, h.ListNotifications)
	g.GET("/unread", h.ParseUserID, h.CountUnRead)
}

// ListNotifications godoc
// @Summary      Get list of notification from latest to oldest
// @Tags         Notification
// @Success      200  {array}	datastore.Notification "List of notifications"
// @Failure      401  {object}  server.ErrorResponse   "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /notification [get]
func (h NotificationHandler) ListNotifications(c *gin.Context) {
	patientID := h.GetUserID(c)
	notifications, err := h.notificationDataStore.ListLatest(patientID)
	if err != nil {
		h.InternalServerError(c, err, "h.notificationDataStore.ListLatest error")
		return
	}
	c.JSON(http.StatusOK, notifications)
}

type CountUnReadNotificationResponse struct {
	Count int `json:"count"`
}

// CountUnRead godoc
// @Summary      Get count of unread notifications
// @Tags         Notification
// @Success      200  {object}	CountUnReadNotificationResponse "Count of the unread notifications"
// @Failure      401  {object}  server.ErrorResponse   "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /notification/unread [get]
func (h NotificationHandler) CountUnRead(c *gin.Context) {
	patientID := h.GetUserID(c)
	count, err := h.notificationDataStore.CountUnRead(patientID)
	if err != nil {
		h.InternalServerError(c, err, "h.notificationDataStore.CountUnRead error")
		return
	}
	c.JSON(http.StatusOK, &CountUnReadNotificationResponse{Count: count})
}
