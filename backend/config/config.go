package config

import (
	"os"
	"strconv"
	"time"
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
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	// AI Model configuration
	aiAPIKey := os.Getenv("GEMINI_API_KEY")
	aiModel := os.Getenv("AI_MODEL")
	if aiModel == "" {
		aiModel = "gemini-2.5-flash"
	}

	aiTimeout := 30 // default 30 seconds
	if timeoutStr := os.Getenv("AI_TIMEOUT"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			aiTimeout = t
		}
	}

	aiMaxRetry := 3 // default 3 retries
	if retryStr := os.Getenv("AI_MAX_RETRY"); retryStr != "" {
		if r, err := strconv.Atoi(retryStr); err == nil {
			aiMaxRetry = r
		}
	}

	return &Config{
		ServerAddress:      serverAddr,
		JWTSecret:          jwtSecret,
		JWTExpirationHours: 24,
		RateLimitRequests:  100,
		RateLimitDuration:  time.Minute,
		AIAPIKey:           aiAPIKey,
		AIModel:            aiModel,
		AITimeout:          aiTimeout,
		AIMaxRetry:         aiMaxRetry,
	}, nil
}
