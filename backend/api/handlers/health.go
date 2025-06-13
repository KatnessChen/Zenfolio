package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetHealthCheck returns a simple health status to verify the API is running
func GetHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}
