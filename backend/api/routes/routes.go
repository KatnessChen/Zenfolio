package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/api/handlers"
	"github.com/transaction-tracker/backend/api/middlewares"
	"github.com/transaction-tracker/backend/config"
)

// SetupRouter configures the API routes
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Create a rate limiter
	rateLimiter := middlewares.NewClientRateLimiter(cfg)

	// Public routes
	r.GET("/health", handlers.GetHealthCheck)
	
	// Login route with rate limiting
	r.POST("/api/v1/login", middlewares.RateLimitMiddleware(rateLimiter), handlers.Login(cfg))

	// API routes (authenticated)
	api := r.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware(cfg))
	api.Use(middlewares.RateLimitMiddleware(rateLimiter))
	{
		// Hello world route for testing authentication
		api.GET("/hello-world", handlers.HelloWorld)
	}
	
	return r
}
