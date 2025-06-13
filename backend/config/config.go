package config

import (
	"os"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	ServerAddress      string
	JWTSecret          string
	JWTExpirationHours int
	RateLimitRequests  int
	RateLimitDuration  time.Duration
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

	return &Config{
		ServerAddress:      serverAddr,
		JWTSecret:          jwtSecret,
		JWTExpirationHours: 24,
		RateLimitRequests:  100,
		RateLimitDuration:  time.Minute,
	}, nil
}
