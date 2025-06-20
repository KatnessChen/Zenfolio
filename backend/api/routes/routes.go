package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/api/handlers"
	"github.com/transaction-tracker/backend/api/middlewares"
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/constants"
	"github.com/transaction-tracker/backend/internal/database"
)

// SetupRouter configures the API routes
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	rateLimiter := middlewares.NewClientRateLimiter(cfg)

	// Initialize database for auth
	dm, err := database.Initialize(cfg)
	if err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	authHandler := handlers.NewAuthHandler(dm.GetDB(), cfg)

	// Public routes (no authentication required)
	publicApi := r.Group(constants.APIVersion)
	publicApi.Use(middlewares.RateLimitMiddleware(rateLimiter))
	{
		publicApi.GET(constants.HealthEndpoint, handlers.GetHealthCheck)
		publicApi.POST(constants.LoginEndpoint, authHandler.Login)
	}

	// Protected API routes
	api := r.Group(constants.APIVersion)
	api.Use(middlewares.AuthMiddleware(dm.GetDB(), cfg))
	api.Use(middlewares.RateLimitMiddleware(rateLimiter))
	{
		// TODO: only allowed admin users to access these routes
		api.GET(constants.HelloWorldEndpoint, handlers.HelloWorld)
		api.GET(constants.DatabaseHealthEndpoint, handlers.DatabaseHealthHandler)
		api.POST(constants.LogoutAllEndpoint, authHandler.LogoutAll)

		api.POST(constants.LogoutEndpoint, authHandler.Logout)
		api.POST(constants.RefreshTokenEndpoint, authHandler.RefreshToken)
		api.GET(constants.MeEndpoint, authHandler.Me)

		api.POST(constants.ExtractTransEndpoint, handlers.ExtractTransactionsHandler(cfg))
	}

	return r
}
