package provider

import (
	"context"
	"fmt"

	"github.com/transaction-tracker/backend/config"
)

// PriceServiceManager manages Price Service integration with error handling
type PriceServiceManager struct {
	client PriceServiceClient
	config *config.Config
}

// NewPriceServiceManager creates a new Price Service manager
func NewPriceServiceManager(cfg *config.Config) *PriceServiceManager {
	client := NewPriceServiceClient(cfg)

	return &PriceServiceManager{
		client: client,
		config: cfg,
	}
}

// GetCurrentPrice retrieves the current price for a single symbol
func (psm *PriceServiceManager) GetCurrentPrice(ctx context.Context, symbol string) (*SymbolCurrentPrice, error) {
	// Fetch from Price Service
	prices, err := psm.client.GetCurrentPrices(ctx, []string{symbol})
	if err != nil {
		return nil, fmt.Errorf("failed to get current price for %s: %w", symbol, err)
	}

	if len(prices) == 0 {
		return nil, fmt.Errorf("no price data returned for symbol %s", symbol)
	}

	return &prices[0], nil
}

// GetCurrentPrices retrieves current prices for multiple symbols
func (psm *PriceServiceManager) GetCurrentPrices(ctx context.Context, symbols []string) ([]SymbolCurrentPrice, error) {
	if len(symbols) == 0 {
		return []SymbolCurrentPrice{}, nil
	}

	// Fetch all prices from Price Service
	prices, err := psm.client.GetCurrentPrices(ctx, symbols)
	if err != nil {
		return nil, fmt.Errorf("failed to get current prices: %w", err)
	}

	return prices, nil
}

// GetHistoricalPrices retrieves historical prices
func (psm *PriceServiceManager) GetHistoricalPrices(ctx context.Context, symbols []string, resolution Resolution, fromDate, toDate string) ([]SymbolHistoricalPrice, error) {
	return psm.client.GetHistoricalPrices(ctx, symbols, resolution, fromDate, toDate)
}

// GetHistoricalPriceAtDate retrieves historical price for a single symbol at a specific date
func (psm *PriceServiceManager) GetHistoricalPriceAtDate(ctx context.Context, symbol string, date string) (*SymbolHistoricalPrice, error) {
	return psm.client.GetHistoricalPriceAtDate(ctx, symbol, date)
}

// HealthCheck performs a health check on the Price Service
func (psm *PriceServiceManager) HealthCheck(ctx context.Context) (*HealthResponse, error) {
	return psm.client.HealthCheck(ctx)
}

// IsServiceHealthy returns true if the Price Service is healthy
func (psm *PriceServiceManager) IsServiceHealthy() bool {
	return psm.client.IsHealthy()
}
