package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/constants"
	"github.com/transaction-tracker/backend/internal/repositories"
	"github.com/transaction-tracker/backend/internal/services"
	"gorm.io/gorm"
)

// AuthMiddleware returns a middleware for JWT authentication
func AuthMiddleware(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	jwtRepo := repositories.NewJWTRepository(db)
	jwtService := services.NewJWTService(cfg, jwtRepo)

	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader(constants.AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": constants.ErrMsgAuthHeaderRequired,
			})
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != constants.BearerTokenPrefix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": constants.ErrMsgInvalidAuthFormat,
			})
			return
		}
		tokenString := parts[1]

		// Special handling for tokens that are only whitespace
		if strings.TrimSpace(tokenString) == "" && tokenString != "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": constants.ErrMsgInvalidAuthFormat,
			})
			return
		}

		// Validate token using the JWT service
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": constants.ErrMsgInvalidToken,
			})
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("token_id", claims.TokenID)

		c.Next()
	}
}
