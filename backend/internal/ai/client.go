package ai

import (
	"context"
	"io"
)

// TransactionData represents extracted transaction information
type TransactionData struct {
	Ticker      string  `json:"ticker"`
	Exchange    string  `json:"exchange"`
	Currency    string  `json:"currency"`
	TradeDate   string  `json:"trade_date"`
	TradeType   string  `json:"trade_type"` // "Buy", "Sell", "Dividends"
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	TradeAmount float64 `json:"trade_amount"`
}

// ImageInput represents an image file for processing
type ImageInput struct {
	Data     io.Reader
	Filename string
	MimeType string
}

// ExtractResponse represents the response from AI model
type ExtractResponse struct {
	Transactions []TransactionData `json:"transactions"`
	Success      bool              `json:"success"`
	Message      string            `json:"message"`
}

// Client defines the interface for AI model interactions
type Client interface {
	// ExtractTransactions processes images and extracts transaction data
	ExtractTransactions(ctx context.Context, images []ImageInput) (*ExtractResponse, error)
	
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
