package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var (
	ErrInvalidNotificationID = server.NewErrorResponse("Invalid notification id")
	ErrNotificationNotFound  = server.NewErrorResponse("Notification not found")
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
	g := r.Group("/notification", h.ParseUserID)
	g.GET("", h.ListNotifications)
	g.PATCH("", h.ReadAll)
	g.GET("/unread", h.CountUnRead)
	g.PATCH("/:id", h.AuthorizedPatientToNotification, h.Read)
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

func (h NotificationHandler) ReadAll(c *gin.Context) {
	patientID := h.GetUserID(c)
	if err := h.notificationDataStore.SetAllAsRead(patientID); err != nil {
		h.InternalServerError(c, err, "h.notificationDataStore.SetAllAsRead error")
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (h NotificationHandler) AuthorizedPatientToNotification(c *gin.Context) {
	patientID := h.GetUserID(c)
	notificationIDStr := c.Param("id")
	if notificationIDStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidNotificationID)
		return
	}
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidNotificationID)
		return
	}

	notification, err := h.notificationDataStore.FindByID(uint(notificationID))
	if err != nil {
		h.InternalServerError(c, err, "h.notificationDataStore.FindByID error")
		return
	}
	if notification == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrNotificationNotFound)
		return
	}
	if notification.PatientID != patientID {
		c.AbortWithStatusJSON(http.StatusForbidden, ErrForbidden)
		return
	}
	c.Set("Notification", notification)
}

func (h NotificationHandler) Read(c *gin.Context) {
	rawNotification, _ := c.Get("Notification")
	notification, _ := rawNotification.(*datastore.Notification)
	if err := h.notificationDataStore.SetAsRead(notification.ID); err != nil {
		h.InternalServerError(c, err, "h.notificationDataStore.SetAsRead error")
		return
	}
	c.AbortWithStatus(http.StatusOK)
}
