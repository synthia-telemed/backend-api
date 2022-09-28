package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/clock"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/id"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/server/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrDoctorNotFound              = server.NewErrorResponse("Doctor not found")
	ErrInitNonScheduledAppointment = server.NewErrorResponse("Cannot initiate room for completed or cancelled appointment")
	ErrDoctorInAnotherRoom         = server.NewErrorResponse("Doctor is in another room. Please close the room before starting a new one")
	ErrNotTimeYet                  = server.NewErrorResponse("The appointment can start 10 minutes early and not later than 3 hours")
	ErrAppointmentIDMissing        = server.NewErrorResponse("Appointment ID is missing")
	ErrAppointmentIDInvalid        = server.NewErrorResponse("Invalid appointment ID")
	ErrAppointmentNotFound         = server.NewErrorResponse("Appointment not found")
	ErrForbidden                   = server.NewErrorResponse("Forbidden")
	ErrDoctorNotInRoom             = server.NewErrorResponse("Doctor isn't currently in any room")
)

type AppointmentHandler struct {
	appointmentDataStore datastore.AppointmentDataStore
	patientDataStore     datastore.PatientDataStore
	doctorDataStore      datastore.DoctorDataStore
	hospitalClient       hospital.SystemClient
	cacheClient          cache.Client
	clock                clock.Clock
	idGenerator          id.Generator
	logger               *zap.SugaredLogger
	server.GinHandler
}

func NewAppointmentHandler(ads datastore.AppointmentDataStore, pds datastore.PatientDataStore, dds datastore.DoctorDataStore, hos hospital.SystemClient, cache cache.Client, clock clock.Clock, id id.Generator, logger *zap.SugaredLogger) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentDataStore: ads,
		patientDataStore:     pds,
		doctorDataStore:      dds,
		hospitalClient:       hos,
		cacheClient:          cache,
		clock:                clock,
		idGenerator:          id,
		logger:               logger,
		GinHandler: server.GinHandler{
			Logger: logger,
		},
	}
}

func (h AppointmentHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/appointment")
	g.GET("/today", middleware.ParseUserID, h.ParseDoctor, h.TodayAppointment)
	g.POST("/:appointmentID", middleware.ParseUserID, h.ParseDoctor, h.AuthorizedDoctorToAppointment, h.InitAppointmentRoom)
	g.POST("/complete", middleware.ParseUserID, h.ParseDoctor, h.CompleteAppointment)
}

type InitAppointmentRoomResponse struct {
	RoomID string `json:"room_id"`
}

// TodayAppointment godoc
// @Summary      Get list of today appointment
// @Tags         Appointment
// @Success      200  {object}	hospital.CategorizedAppointment "List of appointment group by status"
// @Failure      400  {object}  server.ErrorResponse   "Doctor not found"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /appointment/today [get]
func (h AppointmentHandler) TodayAppointment(c *gin.Context) {
	rawDoc, _ := c.Get("Doctor")
	doctor := rawDoc.(*datastore.Doctor)
	loc, _ := time.LoadLocation("Asia/Bangkok")
	appointments, err := h.hospitalClient.ListAppointmentsByDoctorID(context.Background(), doctor.RefID, h.clock.Now().In(loc))
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalClient.ListAppointmentsByDoctorID error")
		return
	}
	c.JSON(http.StatusOK, h.hospitalClient.CategorizeAppointmentByStatus(appointments))
}

// InitAppointmentRoom godoc
// @Summary      Init the appointment room
// @Tags         Appointment
// @Param  		 appointmentID 	path	 integer	true "ID of the appointment"
// @Success      201  {object}  InitAppointmentRoomResponse  "Room ID is return to be used with socket server"
// @Failure      400  {object}  server.ErrorResponse   "Doctor not found"
// @Failure      400  {object}  server.ErrorResponse   "Appointment ID is missing"
// @Failure      400  {object}  server.ErrorResponse   "Invalid appointment ID"
// @Failure      400  {object}  server.ErrorResponse   "Cannot initiate room for completed or cancelled appointment"
// @Failure      400  {object}  server.ErrorResponse   "The appointment can start 10 minutes early and not later than 3 hours"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      403  {object}  server.ErrorResponse   "Forbidden"
// @Failure      404  {object}  server.ErrorResponse   "Appointment not found"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /appointment/{appointmentID} [post]
func (h AppointmentHandler) InitAppointmentRoom(c *gin.Context) {
	rawDoc, _ := c.Get("Doctor")
	doctor := rawDoc.(*datastore.Doctor)

	rawApp, _ := c.Get("Appointment")
	appointment := rawApp.(*hospital.Appointment)
	if appointment.Status != hospital.AppointmentStatusScheduled {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInitNonScheduledAppointment)
		return
	}
	now := h.clock.Now()
	diff := appointment.StartDateTime.Sub(now)
	if diff < -time.Hour*3 || diff > time.Minute*10 {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrNotTimeYet)
		return
	}

	ctx := context.Background()
	// Check the current room that doctor is in
	currentAppID, err := h.cacheClient.Get(ctx, cache.CurrentDoctorAppointmentIDKey(doctor.ID), false)
	if err != nil {
		h.InternalServerError(c, err, "h.cacheClient.Get error")
		return
	}
	if currentAppID != "" {
		if currentAppID != appointment.Id {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrDoctorInAnotherRoom)
			return
		}
		roomID, err := h.cacheClient.Get(ctx, cache.AppointmentRoomIDKey(appointment.Id), false)
		if err != nil {
			h.InternalServerError(c, err, "h.cacheClient.Get error")
			return
		}
		c.JSON(http.StatusCreated, &InitAppointmentRoomResponse{RoomID: roomID})
		return
	}
	// Doctor is not in any room
	roomID, err := h.idGenerator.GenerateRoomID()
	if err != nil {
		h.InternalServerError(c, err, "h.idGenerator.GenerateRoomID error")
		return
	}
	// Set appointment ID that doctor is currently in and room ID of the appointment
	kv := map[string]string{
		cache.CurrentDoctorAppointmentIDKey(doctor.ID): appointment.Id,
		cache.AppointmentRoomIDKey(appointment.Id):     roomID,
	}
	if err := h.cacheClient.MultipleSet(ctx, kv); err != nil {
		h.InternalServerError(c, err, "h.cacheClient.MultipleSet error")
		return
	}

	patient, err := h.patientDataStore.FindByRefID(appointment.PatientID)
	if err != nil {
		h.InternalServerError(c, err, "h.patientDataStore.FindByRefID error")
		return
	}
	info := map[string]string{
		"PatientID":     fmt.Sprintf("%d", patient.ID),
		"DoctorID":      fmt.Sprintf("%d", doctor.ID),
		"AppointmentID": appointment.Id,
		"StartedAt":     now.Format(time.RFC3339),
	}
	if err := h.cacheClient.HashSet(ctx, cache.RoomInfoKey(roomID), info); err != nil {
		h.InternalServerError(c, err, "h.cacheClient.HashSet error")
		return
	}

	// TODO: Push notification to patient
	c.JSON(http.StatusCreated, &InitAppointmentRoomResponse{RoomID: roomID})
}

type CompleteAppointmentRequest struct {
	Status hospital.SettableAppointmentStatus `json:"status" binding:"required,enum" enums:"CANCELLED,COMPLETED"`
}

// CompleteAppointment godoc
// @Summary      Finish the appointment and close the room
// @Tags         Appointment
// @Param 	  	 CompleteAppointmentRequest body CompleteAppointmentRequest true "Status of the appointment"
// @Success      201  ""
// @Failure      400  {object}  server.ErrorResponse   "Doctor not found"
// @Failure      400  {object}  server.ErrorResponse   "Invalid request body"
// @Failure      400  {object}  server.ErrorResponse   "Doctor isn't currently in any room"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /appointment/complete [post]
func (h AppointmentHandler) CompleteAppointment(c *gin.Context) {
	var req CompleteAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}
	rawDoc, _ := c.Get("Doctor")
	doctor := rawDoc.(*datastore.Doctor)

	ctx := context.Background()
	appointmentID, err := h.cacheClient.Get(ctx, cache.CurrentDoctorAppointmentIDKey(doctor.ID), false)
	if err != nil {
		h.InternalServerError(c, err, "h.cacheClient.Get error")
		return
	}
	if appointmentID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrDoctorNotInRoom)
		return
	}
	roomID, err := h.cacheClient.Get(ctx, cache.AppointmentRoomIDKey(appointmentID), false)
	if err != nil {
		h.InternalServerError(c, err, "h.cacheClient.Get error")
		return
	}
	appIDInt, _ := strconv.ParseInt(appointmentID, 10, 32)
	if req.Status == hospital.SettableAppointmentStatusCompleted {
		startedTimeStr, err := h.cacheClient.HashGet(ctx, cache.RoomInfoKey(roomID), "StartedAt")
		if err != nil {
			h.InternalServerError(c, err, "h.cacheClient.Get error")
			return
		}
		startedTime, err := time.Parse(time.RFC3339, startedTimeStr)
		if err != nil {
			h.InternalServerError(c, err, "time.Parse error")
			return
		}
		duration := h.clock.Now().Sub(startedTime).Round(time.Second)
		appointment := datastore.Appointment{
			RefID:       appointmentID,
			Duration:    duration.Seconds(),
			StartedTime: startedTime.UTC(),
		}
		if err := h.appointmentDataStore.Create(&appointment); err != nil {
			h.InternalServerError(c, err, "h.appointmentDataStore.Create error")
			return
		}
	}
	if err := h.cacheClient.Delete(ctx, cache.CurrentDoctorAppointmentIDKey(doctor.ID), cache.AppointmentRoomIDKey(appointmentID), cache.RoomInfoKey(roomID)); err != nil {
		h.InternalServerError(c, err, "h.cacheClient.Delete error")
		return
	}
	if err := h.hospitalClient.SetAppointmentStatus(ctx, int(appIDInt), req.Status); err != nil {
		h.InternalServerError(c, err, "h.hospitalClient.CompleteAppointment error")
		return
	}
	c.AbortWithStatus(http.StatusCreated)
}

func (h AppointmentHandler) ParseDoctor(c *gin.Context) {
	doctorID := h.GetUserID(c)
	doctor, err := h.doctorDataStore.FindByID(doctorID)
	if err != nil {
		h.InternalServerError(c, err, "h.doctorDataStore.FindByID error")
		return
	}
	if doctor == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrDoctorNotFound)
		return
	}
	c.Set("Doctor", doctor)
}

func (h AppointmentHandler) AuthorizedDoctorToAppointment(c *gin.Context) {
	appointmentIDStr := c.Param("appointmentID")
	if appointmentIDStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrAppointmentIDMissing)
		return
	}
	appointmentID, err := strconv.ParseInt(appointmentIDStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrAppointmentIDInvalid)
		return
	}
	apps, err := h.hospitalClient.FindAppointmentByID(context.Background(), int(appointmentID))
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalClient.FindAppointmentByID error")
		return
	}
	if apps == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrAppointmentNotFound)
		return
	}

	rawDoctor, exist := c.Get("Doctor")
	if !exist {
		h.InternalServerError(c, errors.New("c.Get Patient not exist"), "c.Get Doctor not exist")
		return
	}
	doctor, ok := rawDoctor.(*datastore.Doctor)
	if !ok {
		h.InternalServerError(c, errors.New("doctor type casting error"), "Doctor type casting error")
		return
	}
	if apps.Doctor.ID != doctor.RefID {
		c.AbortWithStatusJSON(http.StatusForbidden, ErrForbidden)
		return
	}
	c.Set("Appointment", apps)
}
