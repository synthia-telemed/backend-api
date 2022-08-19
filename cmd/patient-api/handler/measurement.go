package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler/middleware"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
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
	g.POST("/blood-pressure", h.createBloodPressure)
}

type BloodPressureRequest struct {
	DateTime  time.Time `json:"date_time" binding:"required"`
	Systolic  uint      `json:"systolic" binding:"required"`
	Diastolic uint      `json:"diastolic" binding:"required"`
	Pulse     uint      `json:"pulse" binding:"required"`
}

func (h MeasurementHandler) createBloodPressure(c *gin.Context) {
	var req BloodPressureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	bp := &datastore.BloodPressure{
		DateTime:  req.DateTime.UTC(),
		Systolic:  req.Systolic,
		Diastolic: req.Diastolic,
		Pulse:     req.Pulse,
	}
	if err := h.measurementDataStore.CreateBloodPressure(bp); err != nil {
		InternalServerError(c, h.logger, err, "h.measurementDataStore.CreateBloodPressure error")
	}
	c.JSON(http.StatusCreated, bp)
}
