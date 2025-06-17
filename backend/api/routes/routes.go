package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/api/handlers"
	"github.com/transaction-tracker/backend/api/middlewares"
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/constants"
)

// SetupRouter configures the API routes
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	rateLimiter := middlewares.NewClientRateLimiter(cfg)

	r.GET(constants.HealthEndpoint, handlers.GetHealthCheck)
	r.GET("/health/database", handlers.DatabaseHealthHandler)

	r.POST(constants.APIVersion+constants.LoginEndpoint, middlewares.RateLimitMiddleware(rateLimiter), handlers.Login(cfg))

	api := r.Group(constants.APIVersion)
	api.Use(middlewares.AuthMiddleware(cfg))
	api.Use(middlewares.RateLimitMiddleware(rateLimiter))
	{
		api.GET(constants.HelloWorldEndpoint, handlers.HelloWorld)
		api.POST(constants.ExtractTransEndpoint, handlers.ExtractTransactionsHandler(cfg))
		api.GET("/health/database/detailed", handlers.DetailedDatabaseHealthHandler)
	}

	return r
}
