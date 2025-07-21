package provider

import (
	"fmt"
	"strings"

	"github.com/transaction-tracker/price_service/internal/config"
)

// Type represents different stock price providers
type Type string

const (
	AlphaVantage Type = "alpha_vantage"
	// Future providers can be added here
	// ProviderYahoo       Type = "yahoo_finance"
	// ProviderFinnhub     Type = "finnhub"
)

// NewProvider creates a new stock price provider based on configuration
func NewProvider(cfg *config.Config) (StockPriceProvider, error) {
	providerType := Type(strings.ToLower(cfg.StockAPI.Provider))

	switch providerType {
	case AlphaVantage, "":
		if cfg.StockAPI.APIKey == "" {
			return nil, fmt.Errorf("API key is required for Alpha Vantage provider")
		}
		return NewAlphaVantageProvider(cfg.StockAPI.APIKey), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s. Only 'alpha_vantage' is currently supported", providerType)
	}
}
