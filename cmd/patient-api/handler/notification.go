package handler

import (
	"errors"
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
	g.POST("/token", h.ParsePatient, h.SetNotificationToken)
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

// ReadAll godoc
// @Summary      Set all notification as read
// @Tags         Notification
// @Success      200
// @Failure      401  {object}  server.ErrorResponse   "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /notification [patch]
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

// Read godoc
// @Summary      Set specific notification to read
// @Tags         Notification
// @Param  		 notificationID 	path	 integer 	true "ID of the notification"
// @Success      200
// @Failure      400  {object}  server.ErrorResponse   "Invalid notification id"
// @Failure      401  {object}  server.ErrorResponse   "Unauthorized"
// @Failure      403  {object}  server.ErrorResponse   "Patient doesn't own the notification"
// @Failure      404  {object}  server.ErrorResponse   "Notification not found"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /notification/{notificationID} [patch]
func (h NotificationHandler) Read(c *gin.Context) {
	rawNotification, _ := c.Get("Notification")
	notification, _ := rawNotification.(*datastore.Notification)
	if err := h.notificationDataStore.SetAsRead(notification.ID); err != nil {
		h.InternalServerError(c, err, "h.notificationDataStore.SetAsRead error")
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

type SetNotificationTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// SetNotificationToken godoc
// @Summary      Save patient device notification token
// @Tags         Notification
// @Param  		 SetNotificationTokenRequest body SetNotificationTokenRequest true "Notification token"
// @Success      200
// @Failure      400  {object}  server.ErrorResponse   "Invalid request body"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /notification/token [post]
func (h NotificationHandler) SetNotificationToken(c *gin.Context) {
	var req SetNotificationTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	rawPatient, exist := c.Get("Patient")
	if !exist {
		h.InternalServerError(c, errors.New("c.Get Patient not exist"), "c.Get Patient not exist")
		return
	}
	patient, ok := rawPatient.(*datastore.Patient)
	if !ok {
		h.InternalServerError(c, errors.New("patient type casting error"), "Patient type casting error")
		return
	}

	patient.NotificationToken = req.Token
	if err := h.patientDataStore.Save(patient); err != nil {
		h.InternalServerError(c, err, "h.patientDataStore.Save error")
		return
	}
	c.AbortWithStatus(http.StatusOK)
}
