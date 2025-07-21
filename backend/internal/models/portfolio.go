package models

import "time"

// SingleHolding represents basic information about a stock holding
type SingleHolding struct {
	Symbol               string  `json:"symbol"`
	TotalQuantity        float64 `json:"total_quantity"`
	TotalCost            float64 `json:"total_cost"`
	UnitCost             float64 `json:"unit_cost"`
	CurrentPrice         float64 `json:"current_price"`
	MarketValue          float64 `json:"market_value"`
	SimpleReturnRate     float64 `json:"simple_return_rate"`
	AnnualizedReturnRate float64 `json:"annualized_return_rate"`
	RealizedGainLoss     float64 `json:"realized_gain_loss"`
	UnrealizedGainLoss   float64 `json:"unrealized_gain_loss"`
}

// SingleHoldingResponse represents the response structure for stock basic info
type SingleHoldingResponse struct {
	Data      SingleHolding `json:"data"`
	Timestamp time.Time     `json:"timestamp"`
}

// PortfolioAnalysisType represents the type of analysis requested
type PortfolioAnalysisType string

const (
	AnalysisTypeBasic    PortfolioAnalysisType = "basic"
	AnalysisTypeDetailed PortfolioAnalysisType = "detailed"
)
