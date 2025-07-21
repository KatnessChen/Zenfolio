package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/transaction-tracker/backend/config"
)

// PriceServiceClient defines the interface for Price Service operations
type PriceServiceClient interface {
	GetCurrentPrices(ctx context.Context, symbols []string) ([]SymbolCurrentPrice, error)
	GetHistoricalPrices(ctx context.Context, symbols []string, resolution Resolution, fromDate, toDate string) ([]SymbolHistoricalPrice, error)
	HealthCheck(ctx context.Context) (*HealthResponse, error)
	IsHealthy() bool
}

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	CircuitBreakerClosed CircuitBreakerState = iota
	CircuitBreakerOpen
	CircuitBreakerHalfOpen
)

// CircuitBreaker implements a simple circuit breaker pattern
type CircuitBreaker struct {
	maxFailures     int
	resetTimeout    time.Duration
	failureCount    int
	lastFailureTime time.Time
	state           CircuitBreakerState
	mutex           sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        CircuitBreakerClosed,
	}
}

// Execute executes a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Check if we should attempt to reset the circuit breaker
	if cb.state == CircuitBreakerOpen {
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = CircuitBreakerHalfOpen
			cb.failureCount = 0
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	}

	// Execute the function
	err := fn()
	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		if cb.failureCount >= cb.maxFailures {
			cb.state = CircuitBreakerOpen
		}
		return err
	}

	// Success - reset failure count and close circuit if it was half-open
	cb.failureCount = 0
	cb.state = CircuitBreakerClosed
	return nil
}

// priceServiceClient implements PriceServiceClient
type priceServiceClient struct {
	config         *config.Config
	httpClient     *http.Client
	circuitBreaker *CircuitBreaker
	baseURL        string
	apiKey         string
	lastHealthy    time.Time
	mutex          sync.RWMutex
}

// NewPriceServiceClient creates a new Price Service client
func NewPriceServiceClient(cfg *config.Config) PriceServiceClient {
	return &priceServiceClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.PriceService.Timeout,
		},
		circuitBreaker: NewCircuitBreaker(5, 1*time.Minute), // 5 failures, 1 minute reset
		baseURL:        cfg.PriceService.BaseURL,
		apiKey:         cfg.PriceService.APIKey,
		lastHealthy:    time.Now(),
	}
}

// makeRequest makes an HTTP request with retry logic
func (c *priceServiceClient) makeRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.PriceService.MaxRetries; attempt++ {
		var respBody []byte
		err := c.circuitBreaker.Execute(func() error {
			var execErr error
			respBody, execErr = c.doRequest(ctx, method, endpoint, body)
			return execErr
		})

		if err == nil {
			return respBody, nil
		}

		lastErr = err

		// If circuit breaker is open, don't retry
		if strings.Contains(err.Error(), "circuit breaker is open") {
			break
		}

		// Wait before retrying (exponential backoff)
		if attempt < c.config.PriceService.MaxRetries {
			waitTime := time.Duration(attempt+1) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitTime):
			}
		}
	}

	return nil, lastErr
}

// doRequest performs the actual HTTP request
func (c *priceServiceClient) doRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
	url := c.baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	// Log request
	log.Printf("Price Service Request: %s %s", method, url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Price Service Request Failed: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log response
	log.Printf("Price Service Response: %d %s", resp.StatusCode, string(respBody))

	if resp.StatusCode >= 400 {
		var errorResp ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			return nil, fmt.Errorf("price service error: %s - %s", errorResp.Error.Code, errorResp.Error.Message)
		}
		return nil, fmt.Errorf("price service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Update last healthy timestamp on successful response
	c.mutex.Lock()
	c.lastHealthy = time.Now()
	c.mutex.Unlock()

	return respBody, nil
}

// GetCurrentPrices retrieves current prices for the specified symbols
func (c *priceServiceClient) GetCurrentPrices(ctx context.Context, symbols []string) ([]SymbolCurrentPrice, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("symbols list cannot be empty")
	}

	// Build query parameters
	symbolsParam := strings.Join(symbols, ",")
	endpoint := fmt.Sprintf("/api/v1/price/current?symbols=%s", url.QueryEscape(symbolsParam))

	respBody, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get current prices: %w", err)
	}

	var response CurrentPricesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("price service returned unsuccessful response")
	}

	return response.Data, nil
}

// GetHistoricalPrices retrieves historical prices for the specified symbols
func (c *priceServiceClient) GetHistoricalPrices(ctx context.Context, symbols []string, resolution Resolution, fromDate, toDate string) ([]SymbolHistoricalPrice, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("symbols list cannot be empty")
	}

	// Build query parameters
	params := url.Values{}
	params.Set("symbols", strings.Join(symbols, ","))
	params.Set("resolution", string(resolution))
	if fromDate != "" {
		params.Set("from_date", fromDate)
	}
	if toDate != "" {
		params.Set("to_date", toDate)
	}

	endpoint := fmt.Sprintf("/api/v1/price/historical?%s", params.Encode())

	respBody, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical prices: %w", err)
	}

	var response HistoricalPricesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("price service returned unsuccessful response")
	}

	return response.Data, nil
}

// HealthCheck checks the health of the Price Service
func (c *priceServiceClient) HealthCheck(ctx context.Context) (*HealthResponse, error) {
	respBody, err := c.makeRequest(ctx, "GET", "/health", nil)
	if err != nil {
		return nil, fmt.Errorf("health check failed: %w", err)
	}

	var response HealthResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse health response: %w", err)
	}

	return &response, nil
}

// IsHealthy returns true if the service was healthy in the last 5 minutes
func (c *priceServiceClient) IsHealthy() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	const healthCheckTimeout = 5 * time.Minute

	return time.Since(c.lastHealthy) < healthCheckTimeout
}
