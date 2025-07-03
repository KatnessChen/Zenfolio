package ai

import (
	"fmt"
	"log"

	"github.com/transaction-tracker/backend/config"
)

// NewClient creates a new AI client based on the configuration
func NewClient(cfg *config.Config) (Client, error) {
	if cfg.AIAPIKey == "" {
		return nil, fmt.Errorf("AI API key is required")
	}

	aiConfig := &Config{
		APIKey:      cfg.AIAPIKey,
		Model:       cfg.AIModel,
		Timeout:     cfg.AITimeout,
		MaxRetry:    cfg.AIMaxRetry,
		Environment: cfg.Environment,
	}

	log.Printf("Initializing AI client with model: %s", aiConfig.Model)

	client, err := NewAIModelClient(aiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI model client: %w", err)
	}

	return client, nil
}
