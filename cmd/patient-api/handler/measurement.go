package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler/middleware"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type MeasurementHandler struct {
	measurementDataStore datastore.MeasurementDataStore
	logger               *zap.SugaredLogger
}

func NewMeasurementHandler(m datastore.MeasurementDataStore, l *zap.SugaredLogger) *MeasurementHandler {
	return &MeasurementHandler{
		measurementDataStore: m,
		logger:               l,
	}
}

func (h MeasurementHandler) Register(r *gin.RouterGroup) {
	g := r.Group("/measurement", middleware.ParsePatientID)
	g.POST("/blood-pressure", h.CreateBloodPressure)
	g.POST("/glucose", h.CreateGlucose)
}

type BloodPressureRequest struct {
	DateTime  time.Time `json:"date_time" binding:"required"`
	Systolic  uint      `json:"systolic" binding:"required"`
	Diastolic uint      `json:"diastolic" binding:"required"`
	Pulse     uint      `json:"pulse" binding:"required"`
}

// CreateBloodPressure godoc
// @Summary      Record blood pressure
// @Tags         Measurement
// @Param 	  	 BloodPressureRequest body BloodPressureRequest true "Blood pressure information"
// @Success      201  {object}  datastore.BloodPressure
// @Failure      400  {object}  server.ErrorResponse "Invalid request body"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /measurement/blood-pressure [post]
func (h MeasurementHandler) CreateBloodPressure(c *gin.Context) {
	var req BloodPressureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrInvalidRequestBody)
		return
	}
	id, _ := c.Get("patientID")
	patientID := id.(uint)

	bp := &datastore.BloodPressure{
		PatientID: patientID,
		DateTime:  req.DateTime.UTC(),
		Systolic:  req.Systolic,
		Diastolic: req.Diastolic,
		Pulse:     req.Pulse,
	}
	if err := h.measurementDataStore.CreateBloodPressure(bp); err != nil {
		server.InternalServerError(c, h.logger, err, "h.measurementDataStore.CreateBloodPressure error")
	}
	c.JSON(http.StatusCreated, bp)
}

type GlucoseRequest struct {
	DateTime     time.Time `json:"date_time" binding:"required"`
	Value        uint      `json:"value" binding:"required"`
	IsBeforeMeal *bool     `json:"is_before_meal" binding:"required"`
}

// CreateGlucose godoc
// @Summary      Record glucose level
// @Tags         Measurement
// @Param 	  	 GlucoseRequest body GlucoseRequest true "Glucose level information"
// @Success      201  {object}  datastore.Glucose
// @Failure      400  {object}  server.ErrorResponse "Invalid request body"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /measurement/glucose [post]
func (h MeasurementHandler) CreateGlucose(c *gin.Context) {
	var req GlucoseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Debug(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrInvalidRequestBody)
		return
	}
	id, _ := c.Get("patientID")
	patientID := id.(uint)

	g := &datastore.Glucose{
		PatientID:    patientID,
		DateTime:     req.DateTime.UTC(),
		Value:        req.Value,
		IsBeforeMeal: *req.IsBeforeMeal,
	}
	if err := h.measurementDataStore.CreateGlucose(g); err != nil {
		server.InternalServerError(c, h.logger, err, "h.measurementDataStore.CreateGlucose error")
	}
	c.JSON(http.StatusCreated, g)
}
