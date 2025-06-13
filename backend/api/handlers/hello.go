package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloWorld handles the authenticated hello world endpoint
func HelloWorld(c *gin.Context) {
	// Get the user ID from the context (set by the auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
		"user":    userID,
	})
}
