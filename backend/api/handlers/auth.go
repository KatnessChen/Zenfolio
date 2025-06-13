package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/api/middlewares"
	"github.com/transaction-tracker/backend/config"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
}

// Login handles user authentication and returns a JWT token
func Login(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// In a real application, you would validate credentials against a database
		// For this example, we just check if username and password are not empty
		if req.Username == "" || req.Password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// Generate a JWT token
		token, err := middlewares.GenerateToken(req.Username, cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
		})
	}
}
