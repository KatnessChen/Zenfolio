package provider

import (
	"context"

	"github.com/transaction-tracker/price_service/internal/config"
	"github.com/transaction-tracker/price_service/internal/models"
)

// StockPriceProvider defines the interface for stock price data providers
type StockPriceProvider interface {
	// GetCurrentPrices retrieves current prices for multiple symbols
	GetCurrentPrices(ctx context.Context, symbols []string) ([]models.SymbolCurrentPrice, error)

	// GetHistoricalPrices retrieves historical prices for a single symbol
	GetHistoricalPrices(ctx context.Context, symbol string, resolution models.Resolution) (*models.SymbolHistoricalPrice, error)
}

// ThirdPartyProviderMap handles all price-related operations with built-in provider routing
// Implements StockPriceProvider by routing to appropriate third-party providers
type ThirdPartyProviderMap struct {
	alphaVantage *AlphaVantageProvider
	finnhub      *FinnhubProvider
}

// NewThirdPartyProviderMap creates a new price service with both providers
func NewThirdPartyProviderMap(cfg *config.Config) (*ThirdPartyProviderMap, error) {
	// Initialize Alpha Vantage for historical data
	alphaVantage := NewAlphaVantageProvider(cfg.StockAPI.AlphaVantage.APIKey)
	if cfg.StockAPI.AlphaVantage.BaseURL != "" {
		alphaVantage.BaseURL = cfg.StockAPI.AlphaVantage.BaseURL
	}

	// Initialize Finnhub for current prices
	finnhub := NewFinnhubProvider(cfg.StockAPI.Finnhub.APIKey, cfg.StockAPI.Finnhub.BaseURL)

	return &ThirdPartyProviderMap{
		alphaVantage: alphaVantage,
		finnhub:      finnhub,
	}, nil
}

// GetCurrentPrices uses Finnhub (fast, real-time)
func (t *ThirdPartyProviderMap) GetCurrentPrices(ctx context.Context, symbols []string) ([]models.SymbolCurrentPrice, error) {
	return t.finnhub.GetCurrentPrices(ctx, symbols)
}

// GetHistoricalPrices uses Alpha Vantage (comprehensive historical data)
func (t *ThirdPartyProviderMap) GetHistoricalPrices(ctx context.Context, symbol string, resolution models.Resolution) (*models.SymbolHistoricalPrice, error) {
	return t.alphaVantage.GetHistoricalPrices(ctx, symbol, resolution)
}
