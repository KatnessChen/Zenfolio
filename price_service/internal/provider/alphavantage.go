package provider

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/transaction-tracker/price_service/internal/logger"
	"github.com/transaction-tracker/price_service/internal/models"
)

//go:embed testdata/alphavantage-ibm-daily.json
var ibmDailyTestData []byte

// AlphaVantageProvider handles historical stock prices from Alpha Vantage API
type AlphaVantageProvider struct {
	APIKey  string
	BaseURL string
	client  *http.Client
}

func NewAlphaVantageProvider(apiKey string) *AlphaVantageProvider {
	return &AlphaVantageProvider{
		APIKey:  apiKey,
		BaseURL: "https://www.alphavantage.co/query",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (a *AlphaVantageProvider) GetHistoricalPrices(ctx context.Context, symbol string, resolution models.Resolution) (*models.SymbolHistoricalPrice, error) {
	logger.Info("Alpha Vantage GetHistoricalPrices Start", logger.H{
		"symbol":     symbol,
		"resolution": string(resolution),
	})

	params := url.Values{}
	// Map our resolution to Alpha Vantage function
	switch resolution {
	case models.ResolutionDaily:
		params.Set("function", "TIME_SERIES_DAILY")
	case models.ResolutionWeekly:
		params.Set("function", "TIME_SERIES_WEEKLY")
	case models.ResolutionMonthly:
		params.Set("function", "TIME_SERIES_MONTHLY")
	default:
		return nil, fmt.Errorf("unsupported resolution: %s", resolution)
	}

	params.Set("symbol", symbol)
	params.Set("apikey", a.APIKey)
	params.Set("outputsize", "full") // Get full historical data

	var resp []byte
	var err error
	if strings.EqualFold(symbol, "IBM") && resolution == models.ResolutionDaily {
		// Read mock response for IBM daily from file to avoid rate limit
		logger.Info("Using Mock Data for Alpha Vantage", logger.H{
			"symbol":     symbol,
			"resolution": string(resolution),
			"reason":     "rate_limit_avoidance",
		})
		resp, err = readIBMTestData()
		if err != nil {
			logger.Warn("Failed to Read Mock Data", logger.H{
				"symbol": symbol,
				"error":  err.Error(),
			})
			return nil, fmt.Errorf("failed to read IBM mock data: %w", err)
		}
	} else {
		resp, err = a.makeRequest(ctx, params)
		if err != nil {
			return nil, err
		}
	}

	// Parse the response based on resolution
	var timeSeries map[string]map[string]string
	var result map[string]interface{}

	if err := json.Unmarshal(resp, &result); err != nil {
		logger.Warn("Alpha Vantage Response Parse Failed", logger.H{
			"symbol": symbol,
			"error":  err.Error(),
		})
		return nil, fmt.Errorf("failed to parse Alpha Vantage response: %w", err)
	}

	// Check for API error response
	if errorMsg, exists := result["Error Message"]; exists {
		logger.Warn("Alpha Vantage API Error", logger.H{
			"symbol":       symbol,
			"error_message": errorMsg,
		})
		return nil, fmt.Errorf("alpha vantage API error: %v", errorMsg)
	}

	// Check for rate limit response
	if note, exists := result["Note"]; exists {
		logger.Warn("Alpha Vantage Rate Limit", logger.H{
			"symbol": symbol,
			"note":   note,
		})
		return nil, fmt.Errorf("alpha vantage rate limit: %v", note)
	}

	// Extract time series data based on resolution
	switch resolution {
	case models.ResolutionDaily:
		if ts, ok := result["Time Series (Daily)"].(map[string]interface{}); ok {
			timeSeries = make(map[string]map[string]string)
			for date, data := range ts {
				if dataMap, ok := data.(map[string]interface{}); ok {
					convertedData := make(map[string]string)
					for k, v := range dataMap {
						if str, ok := v.(string); ok {
							convertedData[k] = str
						}
					}
					timeSeries[date] = convertedData
				}
			}
		}
	case models.ResolutionWeekly:
		if ts, ok := result["Weekly Time Series"].(map[string]interface{}); ok {
			timeSeries = make(map[string]map[string]string)
			for date, data := range ts {
				if dataMap, ok := data.(map[string]interface{}); ok {
					convertedData := make(map[string]string)
					for k, v := range dataMap {
						if str, ok := v.(string); ok {
							convertedData[k] = str
						}
					}
					timeSeries[date] = convertedData
				}
			}
		}
	case models.ResolutionMonthly:
		if ts, ok := result["Monthly Time Series"].(map[string]interface{}); ok {
			timeSeries = make(map[string]map[string]string)
			for date, data := range ts {
				if dataMap, ok := data.(map[string]interface{}); ok {
					convertedData := make(map[string]string)
					for k, v := range dataMap {
						if str, ok := v.(string); ok {
							convertedData[k] = str
						}
					}
					timeSeries[date] = convertedData
				}
			}
		}
	}

	var prices []models.ClosePrice
	for dateStr, data := range timeSeries {
		// Get closing price
		if closePrice, ok := data["4. close"]; ok {
			if price, err := strconv.ParseFloat(closePrice, 64); err == nil {
				prices = append(prices, models.ClosePrice{
					Date:  dateStr,
					Price: price,
				})
			}
		}
	}

	// Sort prices by date (newest to oldest)
	sort.Slice(prices, func(i, j int) bool {
		d1, err1 := time.Parse("2006-01-02", prices[i].Date)
		d2, err2 := time.Parse("2006-01-02", prices[j].Date)
		if err1 != nil || err2 != nil {
			// Fallback: compare as strings (descending)
			return prices[i].Date > prices[j].Date
		}
		return d1.After(d2)
	})

	logger.Info("Alpha Vantage GetHistoricalPrices Complete", logger.H{
		"symbol":      symbol,
		"resolution":  string(resolution),
		"price_count": len(prices),
	})

	return &models.SymbolHistoricalPrice{
		Symbol:           symbol,
		Resolution:       resolution,
		HistoricalPrices: prices,
	}, nil
}

func (a *AlphaVantageProvider) makeRequest(ctx context.Context, params url.Values) ([]byte, error) {
	reqURL := fmt.Sprintf("%s?%s", a.BaseURL, params.Encode())

	// Log the request URL for debugging (sanitize API key for security)
	sanitizedURL := strings.Replace(reqURL, a.APIKey, "[API_KEY]", -1)
	logger.Info("Alpha Vantage API Request", logger.H{"url": sanitizedURL})

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		logger.Warn("Alpha Vantage API Request Failed", logger.H{
			"url":   sanitizedURL,
			"error": err.Error(),
		})
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Warn("Alpha Vantage API Error Response", logger.H{
			"url":         sanitizedURL,
			"status_code": resp.StatusCode,
		})
		return nil, fmt.Errorf("alpha vantage API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Warn("Alpha Vantage Response Read Failed", logger.H{
			"url":   sanitizedURL,
			"error": err.Error(),
		})
		return nil, err
	}

	logger.Info("Alpha Vantage API Response Success", logger.H{
		"url":           sanitizedURL,
		"status_code":   resp.StatusCode,
		"response_size": len(body),
	})

	return body, nil
}

// readIBMTestData loads the IBM daily mock JSON from the test data directory
func readIBMTestData() ([]byte, error) {
	return ibmDailyTestData, nil
}
