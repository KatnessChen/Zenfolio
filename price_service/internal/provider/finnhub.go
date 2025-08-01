package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/transaction-tracker/price_service/internal/models"
)

// FinnhubProvider handles real-time stock prices from Finnhub API
type FinnhubProvider struct {
	APIKey  string
	BaseURL string
	client  *http.Client
}

// FinnhubQuoteResponse represents the response structure from Finnhub quote API
type FinnhubQuoteResponse struct {
	CurrentPrice  float64 `json:"c"`  // Current price
	Change        float64 `json:"d"`  // Change
	PercentChange float64 `json:"dp"` // Percent change
	HighPrice     float64 `json:"h"`  // High price of the day
	LowPrice      float64 `json:"l"`  // Low price of the day
	OpenPrice     float64 `json:"o"`  // Open price of the day
	PreviousClose float64 `json:"pc"` // Previous close price
	Timestamp     int64   `json:"t"`  // Timestamp
}

func NewFinnhubProvider(apiKey, baseURL string) *FinnhubProvider {
	return &FinnhubProvider{
		APIKey:  apiKey,
		BaseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (f *FinnhubProvider) GetCurrentPrices(ctx context.Context, symbols []string) ([]models.SymbolCurrentPrice, error) {
	var prices []models.SymbolCurrentPrice

	// Finnhub requires individual requests for each symbol for the quote endpoint
	for _, symbol := range symbols {
		price, err := f.getCurrentPriceForSymbol(ctx, symbol)
		if err != nil {
			log.Printf("Error fetching current price for symbol %s: %v", symbol, err)
			// Continue with other symbols instead of failing completely
			continue
		}
		prices = append(prices, price)
	}

	return prices, nil
}

func (f *FinnhubProvider) getCurrentPriceForSymbol(ctx context.Context, symbol string) (models.SymbolCurrentPrice, error) {
	url := fmt.Sprintf("%s/quote?symbol=%s&token=%s", f.BaseURL, symbol, f.APIKey)

	body, err := f.makeRequest(ctx, url)
	if err != nil {
		return models.SymbolCurrentPrice{}, err
	}

	var quoteResp FinnhubQuoteResponse
	if err := json.Unmarshal(body, &quoteResp); err != nil {
		return models.SymbolCurrentPrice{}, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Check if we got valid data (Finnhub returns 0 values for invalid symbols)
	if quoteResp.CurrentPrice == 0 && quoteResp.Timestamp == 0 {
		return models.SymbolCurrentPrice{}, fmt.Errorf("invalid or not found symbol: %s", symbol)
	}

	return models.SymbolCurrentPrice{
		Symbol:        strings.ToUpper(symbol),
		CurrentPrice:  quoteResp.CurrentPrice,
		Change:        quoteResp.Change,
		ChangePercent: quoteResp.PercentChange,
		PreviousClose: quoteResp.PreviousClose,
		Timestamp:     time.Unix(quoteResp.Timestamp, 0),
	}, nil
}

func (f *FinnhubProvider) makeRequest(ctx context.Context, url string) ([]byte, error) {
	// Log the request URL for debugging (sanitize API key for security)
	sanitizedURL := strings.Replace(url, f.APIKey, "[API_KEY]", -1)
	log.Printf("Finnhub API Request: %s", sanitizedURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Finnhub API Error Response: %s", string(body))
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Finnhub API Response Body: %s", string(body))

	return body, nil
}
