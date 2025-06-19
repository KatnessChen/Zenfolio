package types

import (
	"context"
	"io"
)

// TradeType represents the type of financial transaction
type TradeType string

const (
	TradeTypeBuy      TradeType = "buy"
	TradeTypeSell     TradeType = "sell"
	TradeTypeDividend TradeType = "dividend"
)

// TransactionData represents extracted transaction information from AI
// Uses fields that map to the Transaction model structure
type TransactionData struct {
	Symbol          string    `json:"symbol"`           // Maps to Transaction.Symbol
	Type            TradeType `json:"type"`             // Maps to Transaction.Type
	Quantity        float64   `json:"quantity"`         // Maps to Transaction.Quantity
	Price           float64   `json:"price"`            // Maps to Transaction.Price
	Amount          float64   `json:"amount"`           // Maps to Transaction.Amount
	Currency        string    `json:"currency"`         // Maps to Transaction.Currency
	Broker          string    `json:"broker"`           // Maps to Transaction.Broker
	Account         string    `json:"account"`          // Maps to Transaction.Account
	TransactionDate string    `json:"transaction_date"` // Maps to Transaction.TransactionDate (as string for JSON)
	UserNotes       string    `json:"user_notes"`       // Maps to Transaction.UserNotes
	Exchange        string    `json:"exchange"`         // Additional field for AI extraction
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

// AIClient defines the interface for AI model interactions
type AIClient interface {
	// ExtractTransactions processes images and extracts transaction data
	ExtractTransactions(ctx context.Context, images []ImageInput) (*ExtractResponse, error)
	// Close releases any resources held by the client
	Close()
}
