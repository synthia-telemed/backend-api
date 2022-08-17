package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ParsePatientID(c *gin.Context) {
	id := c.Request.Header.Get("X-USER-ID")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Missing patient ID"})
		return
	}
	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid patient ID"})
		return
	}
	c.Set("patientID", uint(uintID))
}
