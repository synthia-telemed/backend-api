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
	"github.com/synthia-telemed/backend-api/pkg/notification"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"go.uber.org/zap"
	"math"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrDoctorNotFound              = server.NewErrorResponse("Doctor not found")
	ErrInitNonScheduledAppointment = server.NewErrorResponse("Cannot join a completed or cancelled appointment")
	ErrDoctorInAnotherRoom         = server.NewErrorResponse("You're in another room. Please close the room before starting a new one")
	ErrNotTimeYet                  = server.NewErrorResponse("The appointment can be started 10 minutes early and not later than 3 hours")
	ErrAppointmentIDMissing        = server.NewErrorResponse("Appointment ID is missing")
	ErrAppointmentIDInvalid        = server.NewErrorResponse("Invalid appointment ID")
	ErrAppointmentNotFound         = server.NewErrorResponse("Appointment not found")
	ErrForbidden                   = server.NewErrorResponse("Forbidden")
	ErrDoctorNotInRoom             = server.NewErrorResponse("You're not currently in any room")
)

type AppointmentHandler struct {
	appointmentDataStore  datastore.AppointmentDataStore
	patientDataStore      datastore.PatientDataStore
	doctorDataStore       datastore.DoctorDataStore
	notificationDataStore datastore.NotificationDataStore
	hospitalClient        hospital.SystemClient
	cacheClient           cache.Client
	clock                 clock.Clock
	idGenerator           id.Generator
	logger                *zap.SugaredLogger
	notificationClient    notification.Client
	server.GinHandler
}

func NewAppointmentHandler(ads datastore.AppointmentDataStore, pds datastore.PatientDataStore, dds datastore.DoctorDataStore, nds datastore.NotificationDataStore, hos hospital.SystemClient, cache cache.Client, clock clock.Clock, id id.Generator, noti notification.Client, logger *zap.SugaredLogger) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentDataStore:  ads,
		patientDataStore:      pds,
		doctorDataStore:       dds,
		notificationDataStore: nds,
		hospitalClient:        hos,
		cacheClient:           cache,
		clock:                 clock,
		idGenerator:           id,
		logger:                logger,
		notificationClient:    noti,
		GinHandler: server.GinHandler{
			Logger: logger,
		},
	}
}

func (h AppointmentHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/appointment", h.ParseUserID, h.ParseDoctor)
	g.GET("", h.ListAppointments)
	g.GET("/:appointmentID", h.AuthorizedDoctorToAppointment, h.GetDoctorAppointmentDetail)
	g.POST("/:appointmentID", h.AuthorizedDoctorToAppointment, h.CanJoinAppointment, h.InitAppointmentRoom, h.SendAppointmentPushNotification)
	g.GET("/:appointmentID/can-join", h.AuthorizedDoctorToAppointment, h.CanJoinAppointment)
	g.POST("/complete", h.CompleteAppointment)
}

type InitAppointmentRoomResponse struct {
	RoomID string `json:"room_id"`
}

type ListAppointmentsRequest struct {
	hospital.ListAppointmentsFilters
	PageNumber int `json:"page_number" form:"page_number" binding:"required"`
	PerPage    int `json:"per_page" form:"per_page" binding:"required"`
}

type ListAppointmentsResponse struct {
	PageNumber   int                             `json:"page_number"`
	PerPage      int                             `json:"per_page"`
	TotalPage    int                             `json:"total_page"`
	TotalItem    int                             `json:"total_item"`
	Appointments []*hospital.AppointmentOverview `json:"appointments"`
}

// ListAppointments godoc
// @Summary      Get list of the appointments with filter
// @Tags         Appointment
// @Param 	  	 ListAppointmentsRequest query ListAppointmentsRequest true "Filter with pagination options for querying"
// @Success      200  {array}	ListAppointmentsResponse "List of appointment overview details with pagination information"
// @Failure      400  {object}  server.ErrorResponse   "Doctor not found"
// @Failure      401  {object}  server.ErrorResponse   "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /appointment [get]
func (h AppointmentHandler) ListAppointments(c *gin.Context) {
	rawDoc, _ := c.Get("Doctor")
	doctor := rawDoc.(*datastore.Doctor)
	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}
	req.DoctorID = &doctor.RefID
	skip := (req.PageNumber - 1) * req.PerPage
	appointments, err := h.hospitalClient.ListAppointmentsWithFilters(context.Background(), &req.ListAppointmentsFilters, req.PerPage, skip)
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalClient.ListAppointmentsByDoctorID error")
		return
	}
	count, err := h.hospitalClient.CountAppointmentsWithFilters(context.Background(), &req.ListAppointmentsFilters)
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalClient.CountAppointmentsWithFilters error")
		return
	}
	res := &ListAppointmentsResponse{
		PageNumber:   req.PageNumber,
		PerPage:      req.PerPage,
		TotalPage:    int(math.Ceil(float64(count) / float64(req.PerPage))),
		TotalItem:    count,
		Appointments: appointments,
	}
	c.JSON(http.StatusOK, res)
}

// GetDoctorAppointmentDetail godoc
// @Summary      Get appointment detail
// @Tags         Appointment
// @Param  		 appointmentID 	path	 integer	true "ID of the appointment"
// @Success      200  {object}  hospital.DoctorAppointment  "Appointment detail"
// @Failure      400  {object}  server.ErrorResponse   "Doctor not found"
// @Failure      400  {object}  server.ErrorResponse   "Appointment ID is missing"
// @Failure      400  {object}  server.ErrorResponse   "Invalid appointment ID"
// @Failure      401  {object}  server.ErrorResponse   "Unauthorized"
// @Failure      403  {object}  server.ErrorResponse   "Forbidden"
// @Failure      404  {object}  server.ErrorResponse   "Appointment not found"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /appointment/{appointmentID} [get]
func (h AppointmentHandler) GetDoctorAppointmentDetail(c *gin.Context) {
	rawApp, _ := c.Get("Appointment")
	appointment := rawApp.(*hospital.DoctorAppointment)
	c.JSON(http.StatusOK, appointment)
}

// CanJoinAppointment godoc
// @Summary      Check if the doctor can join or open the appointment room
// @Tags         Appointment
// @Param  		 appointmentID 	path	 integer	true "ID of the appointment"
// @Success      200
// @Failure      400  {object}  server.ErrorResponse   "Doctor not found"
// @Failure      400  {object}  server.ErrorResponse   "Appointment ID is missing"
// @Failure      400  {object}  server.ErrorResponse   "Invalid appointment ID"
// @Failure      400  {object}  server.ErrorResponse   "Cannot initiate room for completed or cancelled appointment"
// @Failure      400  {object}  server.ErrorResponse   "The appointment can start 10 minutes early and not later than 3 hours"
// @Failure      401  {object}  server.ErrorResponse   "Unauthorized"
// @Failure      403  {object}  server.ErrorResponse   "Forbidden"
// @Failure      404  {object}  server.ErrorResponse   "Appointment not found"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /appointment/{appointmentID}/can-join [get]
func (h AppointmentHandler) CanJoinAppointment(c *gin.Context) {
	rawDoc, _ := c.Get("Doctor")
	doctor := rawDoc.(*datastore.Doctor)

	rawApp, _ := c.Get("Appointment")
	appointment := rawApp.(*hospital.DoctorAppointment)
	if appointment.Status != hospital.AppointmentStatusScheduled {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInitNonScheduledAppointment)
		return
	}

	currentAppID, err := h.cacheClient.Get(context.Background(), cache.CurrentDoctorAppointmentIDKey(doctor.ID), false)
	if err != nil {
		h.InternalServerError(c, err, "h.cacheClient.Get error")
		return
	}
	if currentAppID != "" {
		if currentAppID != appointment.Id {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrDoctorInAnotherRoom)
			return
		}
		roomID, err := h.cacheClient.Get(context.Background(), cache.AppointmentRoomIDKey(appointment.Id), false)
		if err != nil {
			h.InternalServerError(c, err, "h.cacheClient.Get error")
			return
		}
		c.Set("RoomID", roomID)
		return
	}

	now := h.clock.Now()
	diff := appointment.StartDateTime.Sub(now)
	if diff < -time.Hour*3 || diff > time.Minute*10 {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrNotTimeYet)
		return
	}
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
	appointment := rawApp.(*hospital.DoctorAppointment)

	rawRoomID, exist := c.Get("RoomID")
	if exist {
		roomID := rawRoomID.(string)
		c.AbortWithStatusJSON(http.StatusCreated, &InitAppointmentRoomResponse{RoomID: roomID})
		return
	}
	roomID, err := h.idGenerator.GenerateRoomID()
	if err != nil {
		h.InternalServerError(c, err, "h.idGenerator.GenerateRoomID error")
		return
	}

	ctx := context.Background()
	// Set appointment ID that doctor is currently in and room ID of the appointment
	kv := map[string]string{
		cache.CurrentDoctorAppointmentIDKey(doctor.ID): appointment.Id,
		cache.AppointmentRoomIDKey(appointment.Id):     roomID,
	}
	if err := h.cacheClient.MultipleSet(ctx, kv); err != nil {
		h.InternalServerError(c, err, "h.cacheClient.MultipleSet error")
		return
	}

	patient, err := h.patientDataStore.FindByRefID(appointment.Patient.ID)
	if err != nil {
		h.InternalServerError(c, err, "h.patientDataStore.FindByRefID error")
		return
	}
	info := map[string]string{
		"PatientID":     fmt.Sprintf("%d", patient.ID),
		"DoctorID":      fmt.Sprintf("%d", doctor.ID),
		"AppointmentID": appointment.Id,
	}
	if err := h.cacheClient.HashSet(ctx, cache.RoomInfoKey(roomID), info); err != nil {
		h.InternalServerError(c, err, "h.cacheClient.HashSet error")
		return
	}

	c.Set("Patient", patient)
	c.JSON(http.StatusCreated, &InitAppointmentRoomResponse{RoomID: roomID})
}

func (h AppointmentHandler) SendAppointmentPushNotification(c *gin.Context) {
	rawPatient, _ := c.Get("Patient")
	patient, _ := rawPatient.(*datastore.Patient)
	rawApp, _ := c.Get("Appointment")
	appointment := rawApp.(*hospital.DoctorAppointment)

	// Push notification to patient
	notiParam := notification.SendParams{
		ID:    fmt.Sprintf("%d", patient.ID),
		Title: "Your doctor is ready",
		Body:  fmt.Sprintf("%s is ready for the appointment. Tab here to join the room.", appointment.Doctor.FullName),
	}
	notiData := map[string]string{"appointmentID": appointment.Id}
	if err := h.notificationClient.Send(context.Background(), notiParam, notiData); err != nil {
		h.InternalServerErrorWithoutAborting(c, err, "h.notificationClient.Send error")
		return
	}
}

type CompleteAppointmentRequest struct {
	Status hospital.SettableAppointmentStatus `json:"status" binding:"required,enum" enums:"CANCELLED,COMPLETED"`
}

// CompleteAppointment godoc
// @Summary      Finish the appointment and close the room
// @Tags         Appointment
// @Param 	  	 CompleteAppointmentRequest body CompleteAppointmentRequest true "Status of the appointment"
// @Success      201  "Appointment status is set"
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
		durationStr, err := h.cacheClient.HashGet(ctx, cache.RoomInfoKey(roomID), "Duration")
		if err != nil {
			h.InternalServerError(c, err, "h.cacheClient.Get error")
			return
		}
		duration, _ := strconv.ParseInt(durationStr, 10, 32)
		appointment := datastore.Appointment{
			RefID:       appointmentID,
			Duration:    float64(duration),
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
	apps, err := h.hospitalClient.FindDoctorAppointmentByID(context.Background(), int(appointmentID))
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
