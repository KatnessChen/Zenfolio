package routes

import (
	"time"

	"github.com/gin-contrib/cors"
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

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:4173"}, // Common frontend dev ports
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
		publicApi.POST(constants.SignupEndpoint, authHandler.Signup)
	}

	// Protected API routes
	api := r.Group(constants.APIVersion)
	api.Use(middlewares.AuthMiddleware(dm.GetDB(), cfg))
	api.Use(middlewares.RateLimitMiddleware(rateLimiter))
	{
		// TODO: only allowed admin users to access these routes
		api.GET(constants.HelloWorldEndpoint, handlers.HelloWorld)
		api.GET(constants.DatabaseHealthEndpoint, handlers.DatabaseHealthHandler)

		api.POST(constants.LogoutEndpoint, authHandler.Logout)
		api.POST(constants.RefreshTokenEndpoint, authHandler.RefreshToken)
		api.GET(constants.MeEndpoint, authHandler.Me)

		api.POST(constants.ExtractTransEndpoint, handlers.ExtractTransactionsHandler(cfg))
	}

	return r
}
