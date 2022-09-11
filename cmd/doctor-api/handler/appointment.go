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
)

type AppointmentHandler struct {
	patientDataStore datastore.PatientDataStore
	doctorDataStore  datastore.DoctorDataStore
	hospitalClient   hospital.SystemClient
	cacheClient      cache.Client
	clock            clock.Clock
	idGenerator      id.Generator
	logger           *zap.SugaredLogger
	server.GinHandler
}

func NewAppointmentHandler(pds datastore.PatientDataStore, dds datastore.DoctorDataStore, hos hospital.SystemClient, cache cache.Client, clock clock.Clock, id id.Generator, logger *zap.SugaredLogger) *AppointmentHandler {
	return &AppointmentHandler{
		patientDataStore: pds,
		doctorDataStore:  dds,
		hospitalClient:   hos,
		cacheClient:      cache,
		clock:            clock,
		idGenerator:      id,
		logger:           logger,
		GinHandler: server.GinHandler{
			Logger: logger,
		},
	}
}

func (h AppointmentHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/appointment")
	g.POST("/:appointmentID/init", middleware.ParseUserID, h.ParseDoctor, h.AuthorizedDoctorToAppointment, h.InitAppointmentRoom)
}

type InitAppointmentRoomResponse struct {
	RoomID string `json:"room_id"`
}

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
	diff := appointment.DateTime.Sub(now)
	if diff < -time.Hour*3 || diff > time.Minute*10 {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrNotTimeYet)
		return
	}

	ctx := context.Background()
	// Check the current room that doctor is in
	currentAppID, err := h.cacheClient.Get(ctx, CurrentDoctorAppointmentIDKey(doctor.ID), false)
	if err != nil {
		h.InternalServerError(c, err, "h.cacheClient.Get error")
		return
	}
	if currentAppID != "" {
		if currentAppID != appointment.Id {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrDoctorInAnotherRoom)
			return
		}
		roomID, err := h.cacheClient.Get(ctx, AppointmentRoomIDKey(appointment.Id), false)
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
		CurrentDoctorAppointmentIDKey(doctor.ID): appointment.Id,
		AppointmentRoomIDKey(appointment.Id):     roomID,
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
	}
	if err := h.cacheClient.HashSet(ctx, RoomInfoKey(roomID), info); err != nil {
		h.InternalServerError(c, err, "h.cacheClient.HashSet error")
		return
	}

	// TODO: Push notification to patient
	c.JSON(http.StatusCreated, &InitAppointmentRoomResponse{RoomID: roomID})
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

func CurrentDoctorAppointmentIDKey(doctorID uint) string {
	return fmt.Sprintf("doctor:%d:appointment_id", doctorID)
}

func AppointmentRoomIDKey(appointmentID string) string {
	return fmt.Sprintf("appointment:%s:room_id", appointmentID)
}

func RoomInfoKey(roomID string) string {
	return fmt.Sprintf("room:%s", roomID)
}
