package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/models"
	"github.com/transaction-tracker/backend/internal/repositories"
	"github.com/transaction-tracker/backend/internal/services"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related operations
type AuthHandler struct {
	jwtService services.JWTService
	userRepo   repositories.UserRepository
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	jwtRepo := repositories.NewJWTRepository(db)
	userRepo := repositories.NewUserRepository(db)
	jwtService := services.NewJWTService(cfg, jwtRepo)

	return &AuthHandler{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string      `json:"token"`
	User  UserSummary `json:"user"`
}

// UserSummary represents user information returned in responses
type UserSummary struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserInfoResponse represents the response for /me endpoint
type UserInfoResponse struct {
	User         UserSummary `json:"user"`
	ActiveTokens int         `json:"active_tokens"`
}

// Login handles user authentication and returns a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Get device information
	deviceInfo := services.DeviceInfo{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user, deviceInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserSummary{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}

// Logout handles token revocation
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token from header
	tokenString := extractTokenFromHeader(c)
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No token provided",
		})
		return
	}

	// Revoke the token
	if err := h.jwtService.RevokeToken(tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to revoke token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// LogoutAll handles revoking all user tokens
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Revoke all user tokens
	if err := h.jwtService.RevokeAllUserTokens(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to revoke all tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out from all devices",
	})
}

// Me returns current user information
func (h *AuthHandler) Me(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Get user information
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Get active token count
	tokens, err := h.jwtService.GetActiveTokens(userID)
	if err != nil {
		// Log error but don't fail the request
		tokens = []models.JWTToken{}
	}

	c.JSON(http.StatusOK, UserInfoResponse{
		User: UserSummary{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		ActiveTokens: len(tokens),
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get token from header
	tokenString := extractTokenFromHeader(c)
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No token provided",
		})
		return
	}

	// Get device information
	deviceInfo := services.DeviceInfo{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
	}

	// Refresh the token
	newToken, err := h.jwtService.RefreshToken(tokenString, deviceInfo)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

// extractTokenFromHeader extracts the JWT token from the Authorization header
func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Remove "Bearer " prefix
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	return ""
}
