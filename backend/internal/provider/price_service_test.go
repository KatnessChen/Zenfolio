package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/transaction-tracker/backend/config"
)

func TestPriceServiceClient_GetCurrentPrices(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/price/current", r.URL.Path)
		assert.Equal(t, "AAPL,GOOGL", r.URL.Query().Get("symbols"))

		response := CurrentPricesResponse{
			Success: true,
			Data: []SymbolCurrentPrice{
				{
					Symbol:        "AAPL",
					CurrentPrice:  150.00,
					Currency:      "USD",
					Change:        2.50,
					ChangePercent: 1.69,
					PreviousClose: 147.50,
					Timestamp:     time.Now(),
				},
				{
					Symbol:        "GOOGL",
					CurrentPrice:  2800.00,
					Currency:      "USD",
					Change:        -10.00,
					ChangePercent: -0.36,
					PreviousClose: 2810.00,
					Timestamp:     time.Now(),
				},
			},
			Timestamp: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		PriceService: config.PriceServiceConfig{
			BaseURL:            server.URL,
			PriceServiceApiKey: "test-key",
			Timeout:            30 * time.Second,
			MaxRetries:         3,
		},
	}

	client := NewPriceServiceClient(cfg)
	ctx := context.Background()

	prices, err := client.GetCurrentPrices(ctx, []string{"AAPL", "GOOGL"})
	require.NoError(t, err)
	require.Len(t, prices, 2)

	assert.Equal(t, "AAPL", prices[0].Symbol)
	assert.Equal(t, 150.00, prices[0].CurrentPrice)
	assert.Equal(t, "GOOGL", prices[1].Symbol)
	assert.Equal(t, 2800.00, prices[1].CurrentPrice)
}

func TestPriceServiceClient_GetCurrentPrices_Error(t *testing.T) {
	// Mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		response := ErrorResponse{
			Success: false,
			Error: ErrorDetail{
				Code:    ErrInvalidInput,
				Message: "Invalid symbols provided",
			},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		PriceService: config.PriceServiceConfig{
			BaseURL:            server.URL,
			PriceServiceApiKey: "test-key",
			Timeout:            30 * time.Second,
			MaxRetries:         0, // No retries for this test
		},
	}

	client := NewPriceServiceClient(cfg)
	ctx := context.Background()

	_, err := client.GetCurrentPrices(ctx, []string{"INVALID"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "INVALID_INPUT")
}

func TestPriceServiceClient_HealthCheck(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/health", r.URL.Path)

		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Service:   "price-service",
			Version:   "1.0.0",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		PriceService: config.PriceServiceConfig{
			BaseURL:            server.URL,
			PriceServiceApiKey: "test-key",
			Timeout:            30 * time.Second,
			MaxRetries:         3,
		},
	}

	client := NewPriceServiceClient(cfg)
	ctx := context.Background()

	health, err := client.HealthCheck(ctx)
	require.NoError(t, err)
	assert.Equal(t, "healthy", health.Status)
	assert.Equal(t, "price-service", health.Service)
}

func TestPriceServiceManager_GetCurrentPrice(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := CurrentPricesResponse{
			Success: true,
			Data: []SymbolCurrentPrice{
				{
					Symbol:       "AAPL",
					CurrentPrice: 150.00,
					Currency:     "USD",
					Timestamp:    time.Now(),
				},
			},
			Timestamp: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		PriceService: config.PriceServiceConfig{
			BaseURL:            server.URL,
			PriceServiceApiKey: "test-key",
			Timeout:            30 * time.Second,
			MaxRetries:         3,
		},
	}

	manager := NewPriceServiceManager(cfg)
	ctx := context.Background()

	price, err := manager.GetCurrentPrice(ctx, "AAPL")
	require.NoError(t, err)
	assert.Equal(t, "AAPL", price.Symbol)
	assert.Equal(t, 150.00, price.CurrentPrice)
}

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)

	// First failure
	err1 := cb.Execute(func() error {
		return assert.AnError
	})
	assert.Error(t, err1)

	// Second failure - should open circuit
	err2 := cb.Execute(func() error {
		return assert.AnError
	})
	assert.Error(t, err2)

	// Third call - circuit should be open
	err3 := cb.Execute(func() error {
		return nil // This won't be called
	})
	assert.Error(t, err3)
	assert.Contains(t, err3.Error(), "circuit breaker is open")

	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)

	// Circuit should be half-open and allow one call
	err4 := cb.Execute(func() error {
		return nil // Success
	})
	assert.NoError(t, err4)

	// Circuit should be closed now
	err5 := cb.Execute(func() error {
		return nil
	})
	assert.NoError(t, err5)
}
