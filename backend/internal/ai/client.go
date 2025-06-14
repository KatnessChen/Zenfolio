package ai

import (
	"context"
	"io"

	"github.com/transaction-tracker/backend/internal/constants"
)

// TransactionData represents extracted transaction information
type TransactionData struct {
	Ticker      string     `json:"ticker"`
	TickerLabel string     `json:"ticker_label"`
	Exchange    string     `json:"exchange"`
	Currency    string     `json:"currency"`
	TradeDate   string     `json:"trade_date"`
	TradeType   TradeType  `json:"trade_type"`
	Quantity    float64    `json:"quantity"`
	Price       float64    `json:"price"`
	TradeAmount float64    `json:"trade_amount"`
}

// TradeType represents the type of trade transaction
type TradeType string

// Available trade types
const (
	Buy       TradeType = TradeType(constants.TradeTypeBuy)
	Sell      TradeType = TradeType(constants.TradeTypeSell)
	Dividends TradeType = TradeType(constants.TradeTypeDividends)
)

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
