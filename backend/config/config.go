package config

import (
	"os"
	"strconv"
	"strings"
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
	AIAPIKey   string
	AIModel    string
	AITimeout  int
	AIMaxRetry int
}

// Load loads the configuration from environment variables
// with sensible defaults when variables are not set
func Load() (*Config, error) {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = constants.DefaultServerAddr
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	// AI Model configuration
	aiModel := os.Getenv("AI_MODEL")
	if aiModel == "" {
		aiModel = constants.DefaultAIModel
	}

	// Set API key based on AI model
	var aiAPIKey string
	if strings.Contains(aiModel, "gemini") {
		aiAPIKey = os.Getenv("GEMINI_API_KEY")
	} else {
		// For other AI models, you can set different environment variables
		aiAPIKey = os.Getenv("AI_API_KEY")
	}

	aiTimeout := constants.DefaultAITimeout
	if timeoutStr := os.Getenv("AI_TIMEOUT"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			aiTimeout = t
		}
	}

	aiMaxRetry := constants.DefaultAIMaxRetry
	if retryStr := os.Getenv("AI_MAX_RETRY"); retryStr != "" {
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
