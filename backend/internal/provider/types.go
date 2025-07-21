package provider

import "time"

// Resolution types for historical data
type Resolution string

const (
	ResolutionDaily    Resolution = "daily"
	ResolutionWeekly   Resolution = "weekly"
	ResolutionMonthly  Resolution = "monthly"
	ResolutionIntraday Resolution = "intraday"
)

// SymbolCurrentPrice represents current price data for a symbol
type SymbolCurrentPrice struct {
	Symbol        string    `json:"symbol"`
	CurrentPrice  float64   `json:"current_price"`
	Currency      string    `json:"currency"`
	Change        float64   `json:"change"`         // absolute difference (current_price - previous_close)
	ChangePercent float64   `json:"change_percent"` // relative difference (change/previous_close)
	PreviousClose float64   `json:"previous_close"`
	Timestamp     time.Time `json:"timestamp"`
}

// ClosePrice represents a date-price pair
type ClosePrice struct {
	Date  string  `json:"date"` // YYYY-MM-DD format
	Price float64 `json:"price"`
}

// SymbolHistoricalPrice represents historical price data for a symbol
type SymbolHistoricalPrice struct {
	Symbol           string       `json:"symbol"`
	Resolution       Resolution   `json:"resolution"`
	HistoricalPrices []ClosePrice `json:"historical_prices"`
}

// ErrorCode represents error codes from Price Service
type ErrorCode string

const (
	ErrSymbolNotFound     ErrorCode = "SYMBOL_NOT_FOUND"
	ErrMarketClosed       ErrorCode = "MARKET_CLOSED"
	ErrRateLimitExceeded  ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrInvalidInput       ErrorCode = "INVALID_INPUT"
	ErrUnauthorized       ErrorCode = "UNAUTHORIZED"
)

// ErrorResponse represents the standard error response format from Price Service
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// SuccessResponse represents successful responses from Price Service
type SuccessResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// CurrentPricesResponse represents the response from /api/v1/price/current
type CurrentPricesResponse struct {
	Success   bool                 `json:"success"`
	Data      []SymbolCurrentPrice `json:"data"`
	Timestamp time.Time            `json:"timestamp"`
}

// HistoricalPricesResponse represents the response from /api/v1/price/historical
type HistoricalPricesResponse struct {
	Success   bool                    `json:"success"`
	Data      []SymbolHistoricalPrice `json:"data"`
	Timestamp time.Time               `json:"timestamp"`
}

// HealthResponse represents the response from /health endpoint
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// GetCurrentPricesRequest represents the request parameters for getting current prices
type GetCurrentPricesRequest struct {
	Symbols []string `json:"symbols"`
}

// GetHistoricalPricesRequest represents the request parameters for getting historical prices
type GetHistoricalPricesRequest struct {
	Symbols    []string   `json:"symbols"`
	Resolution Resolution `json:"resolution"`
	FromDate   string     `json:"from_date"` // YYYY-MM-DD format
	ToDate     string     `json:"to_date"`   // YYYY-MM-DD format
}
