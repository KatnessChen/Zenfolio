package ai

import (
	"context"

	"github.com/transaction-tracker/backend/internal/types"
)

// Client defines the interface for AI model interactions
type Client interface {
	// ExtractTransactions processes images and extracts transaction data
	ExtractTransactions(ctx context.Context, images []types.ImageInput) (*types.ExtractResponse, error)
	// Health checks if the AI client is working properly
	Health(ctx context.Context) error
	// Close closes the client and cleans up resources
	Close() error
}

// Config holds configuration for AI clients
type Config struct {
	APIKey   string
	Model    string
	Timeout  int // timeout in seconds
	MaxRetry int
}
