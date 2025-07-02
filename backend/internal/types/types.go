package types

import (
	"io"
)

// TradeType represents the type of financial transaction
type TradeType string

const (
	TradeTypeBuy      TradeType = "Buy"
	TradeTypeSell     TradeType = "Sell"
	TradeTypeDividend TradeType = "Dividends"
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
	Exchange        string    `json:"exchange"`         // Maps to Transaction.Exchange
}

// FileInput represents an image file for processing
type FileInput struct {
	Data     io.Reader
	Filename string
	MimeType string
}

// ExtractResponseData represents the data part of extract response
type ExtractResponseData struct {
	Transactions     []TransactionData `json:"transactions"`
	TransactionCount int               `json:"transaction_count"`
	FileName         string            `json:"file_name"`
}

// ExtractResponse represents the response from AI model
type ExtractResponse struct {
	Data    *ExtractResponseData `json:"data,omitempty"`
	Success bool                 `json:"success"`
	Message string               `json:"message"`
}

// TradeTypeFromString converts string to TradeType
func TradeTypeFromString(s string) (TradeType, bool) {
	switch s {
	case "Buy":
		return TradeTypeBuy, true
	case "Sell":
		return TradeTypeSell, true
	case "Dividends":
		return TradeTypeDividend, true
	default:
		return "", false
	}
}
