package config

import (
	"os"
	"strconv"
	"time"

	"github.com/transaction-tracker/backend/internal/constants"
)

// Config holds all configuration for the application
type Config struct {
	ServerAddress      string
	JWTSecret          string
	JWTExpirationHours int
	RateLimitRequests  int
	RateLimitDuration  time.Duration
	// AI Model Configuration
	AIAPIKey     string
	AIModel      string
	AITimeout    int
	AIMaxRetry   int
}

// Load loads the configuration from environment variables
// with sensible defaults when variables are not set
func Load() (*Config, error) {
	serverAddr := os.Getenv(constants.EnvServerAddr)
	if serverAddr == "" {
		serverAddr = constants.DefaultServerAddr
	}

	jwtSecret := os.Getenv(constants.EnvJWTSecret)
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	// AI Model configuration
	aiAPIKey := os.Getenv(constants.EnvGeminiAPIKey)
	aiModel := os.Getenv(constants.EnvAIModel)
	if aiModel == "" {
		aiModel = constants.DefaultAIModel
	}

	aiTimeout := constants.DefaultAITimeout // default 30 seconds
	if timeoutStr := os.Getenv(constants.EnvAITimeout); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			aiTimeout = t
		}
	}

	aiMaxRetry := constants.DefaultAIMaxRetry // default 3 retries
	if retryStr := os.Getenv(constants.EnvAIMaxRetry); retryStr != "" {
		if r, err := strconv.Atoi(retryStr); err == nil {
			aiMaxRetry = r
		}
	}

	return &Config{
		ServerAddress:      serverAddr,
		JWTSecret:          jwtSecret,
		JWTExpirationHours: constants.DefaultJWTExpiry,
		RateLimitRequests:  constants.DefaultRateLimit,
		RateLimitDuration:  time.Minute,
		AIAPIKey:           aiAPIKey,
		AIModel:            aiModel,
		AITimeout:          aiTimeout,
		AIMaxRetry:         aiMaxRetry,
	}, nil
}
