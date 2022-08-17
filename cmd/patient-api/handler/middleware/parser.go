package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"net/http"
	"strconv"
)

var (
	ErrMissingPatientID = handler.ErrorResponse{Message: "Missing patient ID"}
	ErrInvalidPatientID = handler.ErrorResponse{Message: "Invalid patient ID"}
)

func ParsePatientID(c *gin.Context) {
	id := c.Request.Header.Get("X-USER-ID")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrMissingPatientID)
		return
	}
	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidPatientID)
		return
	}
	c.Set("patientID", uint(uintID))
}
