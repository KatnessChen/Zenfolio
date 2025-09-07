package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/transaction-tracker/backend/internal/constants"
)

// Config holds all configuration for the application
type Config struct {
	ServerAddress      string
	Environment        string // "development" or "production"
	JWTSecret          string
	JWTExpirationHours int
	RateLimitRequests  int
	RateLimitDuration  time.Duration
	// AI Model Configuration
	AIAPIKey   string
	AIModel    string
	AITimeout  int
	AIMaxRetry int
	// Price Service Configuration
	PriceService PriceServiceConfig
}

// PriceServiceConfig holds configuration for Price Service integration
type PriceServiceConfig struct {
	BaseURL            string
	PriceServiceApiKey string
	Timeout            time.Duration
	MaxRetries         int
}

// Load loads the configuration from environment variables
// with sensible defaults when variables are not set
func Load() (*Config, error) {
	// Load .env file if it exists
	envPath := ".env"
	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("Warning: .env file not found at %s or error loading it: %v\n", envPath, err)
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = constants.DefaultServerAddr
	}

	environment := getEnvOrDefault("ENVIRONMENT", "development")

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

	jwtExpirationHours := getEnvOrDefaultInt("JWT_EXPIRATION_HOURS", constants.DefaultJWTExpiry)

	// Price Service configuration
	priceServiceConfig := PriceServiceConfig{
		BaseURL:            getEnvOrDefault("PRICE_SERVICE_URL", "http://localhost:8081"),
		PriceServiceApiKey: getEnvOrDefault("PRICE_SERVICE_API_KEY", ""),
		Timeout:            time.Duration(getEnvOrDefaultInt("PRICE_SERVICE_TIMEOUT", 30)) * time.Second,
		MaxRetries:         getEnvOrDefaultInt("PRICE_SERVICE_MAX_RETRIES", 3),
	}

	return &Config{
		ServerAddress:      serverAddr,
		Environment:        environment,
		JWTSecret:          jwtSecret,
		JWTExpirationHours: jwtExpirationHours,
		RateLimitRequests:  constants.DefaultRateLimit,
		RateLimitDuration:  time.Minute,
		AIAPIKey:           aiAPIKey,
		AIModel:            aiModel,
		AITimeout:          aiTimeout,
		AIMaxRetry:         aiMaxRetry,
		PriceService:       priceServiceConfig,
	}, nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvOrDefaultInt returns environment variable as int or default if not set
func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
