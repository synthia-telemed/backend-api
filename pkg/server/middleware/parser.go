package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ParseUserID(c *gin.Context) {
	id := c.Request.Header.Get("X-USER-ID")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing user ID"})
		return
	}
	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid user ID"})
		return
	}
	c.Set("UserID", uint(uintID))
}
