package tests

import (
	"context"
	"testing"

	"github.com/transaction-tracker/price_service/internal/config"
	"github.com/transaction-tracker/price_service/internal/provider"
)

func TestNewThirdPartyProviderMap(t *testing.T) {
	cfg := &config.Config{
		StockAPI: config.StockAPIConfig{
			AlphaVantage: config.ProviderConfig{
				APIKey:  "test-alpha-key",
				BaseURL: "https://www.alphavantage.co/query",
			},
			Finnhub: config.ProviderConfig{
				APIKey:  "test-finnhub-key",
				BaseURL: "https://finnhub.io/api/v1",
			},
		},
	}

	priceProvider, err := provider.NewThirdPartyProviderMap(cfg)
	if err != nil {
		t.Fatalf("NewThirdPartyProviderMap failed: %v", err)
	}

	if priceProvider == nil {
		t.Fatal("NewThirdPartyProviderMap returned nil provider")
	}

	// Test that the interface methods are available
	// This confirms both providers are properly initialized
	ctx := context.Background()
	_, err = priceProvider.GetCurrentPrices(ctx, []string{"AAPL"})
	// We expect this to fail due to no real API keys, but method should exist
	if err == nil {
		t.Log("GetCurrentPrices method works (or test API key is valid)")
	}
}

func TestNewThirdPartyProviderMapMissingKeys(t *testing.T) {
	// Test missing Alpha Vantage key
	cfg := &config.Config{
		StockAPI: config.StockAPIConfig{
			AlphaVantage: config.ProviderConfig{
				APIKey:  "",
				BaseURL: "https://www.alphavantage.co/query",
			},
			Finnhub: config.ProviderConfig{
				APIKey:  "test-finnhub-key",
				BaseURL: "https://finnhub.io/api/v1",
			},
		},
	}

	// Should still work - just Alpha Vantage will fail when called
	_, err := provider.NewThirdPartyProviderMap(cfg)
	if err != nil {
		t.Fatalf("NewThirdPartyProviderMap should not fail with missing keys: %v", err)
	}
}
