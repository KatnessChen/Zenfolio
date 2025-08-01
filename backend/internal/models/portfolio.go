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
	SingleHolding
	Timestamp time.Time `json:"timestamp"`
}

// AllHoldingsResponse represents the response structure for all holdings
type AllHoldingsResponse struct {
	Holdings  []SingleHolding `json:"holdings"`
	Timestamp time.Time       `json:"timestamp"`
}

// PortfolioSummary represents the overall portfolio summary
type PortfolioSummary struct {
	Timestamp             time.Time `json:"timestamp"`
	Currency              string    `json:"currency"`
	MarketValue           float64   `json:"market_value"`
	TotalCost             float64   `json:"total_cost"`
	TotalReturn           float64   `json:"total_return"`
	TotalReturnPercentage float64   `json:"total_return_percentage"`
	HoldingsCount         int       `json:"holdings_count"`
	HasTransactions       bool      `json:"has_transactions"`
	AnnualizedReturnRate  float64   `json:"annualized_return_rate"`
	LastUpdated           time.Time `json:"last_updated"`
}

// PortfolioAnalysisType represents the type of analysis requested
type PortfolioAnalysisType string

const (
	AnalysisTypeBasic    PortfolioAnalysisType = "basic"
	AnalysisTypeDetailed PortfolioAnalysisType = "detailed"
)

// TimeFrame represents supported timeframes for historical data
type TimeFrame string

const (
	TimeFrame1D  TimeFrame = "1D"
	TimeFrame1W  TimeFrame = "1W"
	TimeFrame1M  TimeFrame = "1M"
	TimeFrame3M  TimeFrame = "3M"
	TimeFrame6M  TimeFrame = "6M"
	TimeFrameYTD TimeFrame = "YTD"
	TimeFrame1Y  TimeFrame = "1Y"
	TimeFrame5Y  TimeFrame = "5Y"
	TimeFrameALL TimeFrame = "ALL"
)

// Granularity represents data point frequency
type Granularity string

const (
	GranularityHourly  Granularity = "hourly"
	GranularityDaily   Granularity = "daily"
	GranularityWeekly  Granularity = "weekly"
	GranularityMonthly Granularity = "monthly"
)

// TotalValueDataPoint represents a single data point in the historical chart
type TotalValueDataPoint struct {
	Timestamp        time.Time `json:"timestamp"`
	TotalValue       float64   `json:"market_value"`
	DayChange        float64   `json:"day_change"`
	DayChangePercent float64   `json:"day_change_percent"`
}

// TotalValueTrendSummary represents summary statistics for the time period
type TotalValueTrendSummary struct {
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Volatility    float64 `json:"volatility"`
	MaxValue      float64 `json:"max_value"`
	MinValue      float64 `json:"min_value"`
}

// HistoricalTotalValueResponse represents the complete response for historical data
type HistoricalTotalValueResponse struct {
	TimeFrame   TimeFrame   `json:"timeframe"`
	Granularity Granularity `json:"granularity"`
	Period      struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	} `json:"period"`
	DataPoints []TotalValueDataPoint  `json:"data_points"`
	Summary    TotalValueTrendSummary `json:"summary"`
}
