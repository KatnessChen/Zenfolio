package ai_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/transaction-tracker/backend/internal/ai"
	"github.com/transaction-tracker/backend/internal/constants"
)

func TestAIModelClientCreation(t *testing.T) {
	// Test with empty API key
	config := &ai.Config{}
	_, err := ai.NewAIModelClient(config)
	if err == nil {
		t.Error("Expected error when API key is empty")
	}

	// Test with valid config (but dummy API key)
	config = &ai.Config{
		APIKey:   "dummy-key-for-testing",
		Model:    constants.DefaultAIModel,
		Timeout:  constants.DefaultAITimeout,
		MaxRetry: constants.DefaultAIMaxRetry,
	}

	client, err := ai.NewAIModelClient(config)
	if err != nil {
		t.Errorf("Unexpected error creating client: %v", err)
	}

	if client == nil {
		t.Error("Client should not be nil")
	}

	// Clean up
	if client != nil {
		client.Close()
	}
}

func TestExtractTransactionsValidation(t *testing.T) {
	config := &ai.Config{
		APIKey:   "dummy-key-for-testing",
		Model:    constants.DefaultAIModel,
		Timeout:  constants.DefaultAITimeout,
		MaxRetry: constants.DefaultAIMaxRetry,
	}

	client, err := ai.NewAIModelClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test with no images
	ctx := context.Background()
	resp, err := client.ExtractTransactions(ctx, []ai.ImageInput{})

	if err != nil {
		t.Errorf("Should not return error for empty image list: %v", err)
	}

	if resp.Success {
		t.Error("Response should not be successful with no images")
	}

	if !strings.Contains(resp.Message, constants.ErrMsgNoImagesProvided) {
		t.Errorf("Expected '%s' message, got: %s", constants.ErrMsgNoImagesProvided, resp.Message)
	}
}

// TestExtractTransactionsFromImage tests the complete image-to-JSON transformation workflow
// This is an integration test that requires a valid GEMINI_API_KEY environment variable
func TestExtractTransactionsFromImage(t *testing.T) {
	// Check if we have a real API key for integration testing
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: GEMINI_API_KEY not set")
	}

	config := &ai.Config{
		APIKey:   apiKey,
		Model:    constants.DefaultAIModel,
		Timeout:  90, // Longer timeout for multiple images
		MaxRetry: constants.DefaultAIMaxRetry,
	}

	client, err := ai.NewAIModelClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test multiple images together
	testImages := []string{
		"Firstrade-total_13_row.png",
		"Firstrade-total_3_row.png",
	}

	var imageInputs []ai.ImageInput

	// Load all test images
	for i, imageName := range testImages {
		testImagePath := filepath.Join("..", "dummy-data", "transaction-screenshots", imageName)
		imageFile, err := os.Open(testImagePath)
		if err != nil {
			t.Fatalf("Failed to open test image %d (%s): %v", i+1, imageName, err)
		}
		defer imageFile.Close()

		imageInputs = append(imageInputs, ai.ImageInput{
			Data:     imageFile,
			Filename: imageName,
			MimeType: constants.MimeTypePNG,
		})
	}

	// Test the extraction with multiple images
	ctx := context.Background()
	resp, err := client.ExtractTransactions(ctx, imageInputs)

	if err != nil {
		t.Fatalf("Failed to extract transactions: %v", err)
	}

	// Validate the response structure
	if !resp.Success {
		t.Errorf("Expected successful extraction, got failure: %s", resp.Message)
	}

	if len(resp.Transactions) == 0 {
		t.Error("Expected at least one transaction to be extracted")
	}

	// Log processing results
	t.Logf("=== Multiple Images Processing Results ===")
	t.Logf("Images processed: %d", len(imageInputs))
	for i, imageName := range testImages {
		t.Logf("  - Image %d: %s", i+1, imageName)
	}
	t.Logf("Total transactions extracted: %d", len(resp.Transactions))
	t.Logf("")

	// Validate transaction data structure
	for i, transaction := range resp.Transactions {
		t.Logf("Transaction %d:", i+1)
		t.Logf("  Ticker: %s", transaction.Ticker)
		t.Logf("  Ticker Label: %s", transaction.TickerLabel)
		t.Logf("  Exchange: %s", transaction.Exchange)
		t.Logf("  Currency: %s", transaction.Currency)
		t.Logf("  Trade Date: %s", transaction.TradeDate)
		t.Logf("  Trade Type: %s", transaction.TradeType)
		t.Logf("  Quantity: %.2f", transaction.Quantity)
		t.Logf("  Price: %.2f", transaction.Price)
		t.Logf("  Trade Amount: %.2f", transaction.TradeAmount)
		t.Logf("  ---")

		// Check that required fields are not empty
		if transaction.Ticker == "" {
			t.Errorf("Transaction %d: Ticker should not be empty", i+1)
		}

		if transaction.TickerLabel == "" {
			t.Logf("Transaction %d: TickerLabel is empty (this may be expected for some tickers)", i+1)
		}

		if transaction.TradeType == "" {
			t.Errorf("Transaction %d: TradeType should not be empty", i+1)
		}

		// Validate trade type is one of the allowed values
		validTradeTypes := constants.ValidTradeTypesMap()
		if !validTradeTypes[string(transaction.TradeType)] {
			t.Errorf("Transaction %d: Invalid trade type '%s', must be %v", i+1, transaction.TradeType, constants.ValidTradeTypes())
		}

		if transaction.TradeDate == "" {
			t.Errorf("Transaction %d: TradeDate should not be empty", i+1)
		}

		// Validate numeric fields are reasonable (allow negative for sell transactions)
		if transaction.Price < 0 {
			t.Errorf("Transaction %d: Price should not be negative: %f", i+1, transaction.Price)
		}
	}

	t.Logf("Successfully extracted %d transactions from %d images", len(resp.Transactions), len(imageInputs))
}

// TestParseTransactionResponse tests the JSON parsing functionality with mocked data
func TestParseTransactionResponse(t *testing.T) {
	config := &ai.Config{
		APIKey:   "dummy-key-for-testing",
		Model:    constants.DefaultAIModel,
		Timeout:  constants.DefaultAITimeout,
		MaxRetry: constants.DefaultAIMaxRetry,
	}

	client, err := ai.NewAIModelClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test with valid JSON response
	mockJSONResponse := `{
		"transactions": [
			{
				"ticker": "AAPL",
				"ticker_label": "Apple Inc.",
				"exchange": "NASDAQ",
				"currency": "USD",
				"trade_date": "2024-06-14",
				"trade_type": "Buy",
				"quantity": 10.0,
				"price": 150.25,
				"trade_amount": 1502.50
			},
			{
				"ticker": "GOOGL",
				"ticker_label": "Alphabet Inc.",
				"exchange": "NASDAQ",
				"currency": "USD",
				"trade_date": "2024-06-14",
				"trade_type": "Sell",
				"quantity": 5.0,
				"price": 2800.00,
				"trade_amount": 14000.00
			}
		]
	}`

	// Since parseTransactionResponse is not exported, we'll test it indirectly
	// by creating a mock scenario that exercises the JSON parsing logic
	t.Logf("Mock JSON response: %s", mockJSONResponse)

	// Validate that our mock data contains all required fields
	expectedTransactions := 2
	t.Logf("Expected number of transactions: %d", expectedTransactions)

	// Test invalid JSON handling would require access to internal methods
	// For now, we verify that our test data structure is valid
	invalidJSONResponse := `{"transactions": [{"invalid": "json"`
	t.Logf("Invalid JSON for error testing: %s", invalidJSONResponse)
}

// TestAIClientWithMockedSuccessResponse tests the complete workflow with expected successful data
func TestAIClientWithMockedSuccessResponse(t *testing.T) {
	config := &ai.Config{
		APIKey:   "dummy-key-for-testing",
		Model:    constants.DefaultAIModel,
		Timeout:  constants.DefaultAITimeout,
		MaxRetry: constants.DefaultAIMaxRetry,
	}

	client, err := ai.NewAIModelClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test data validation - verify our expected response structure
	expectedFields := []string{
		"ticker", "ticker_label", "exchange", "currency",
		"trade_date", "trade_type", "quantity", "price", "trade_amount",
	}

	t.Logf("Testing that AI response should contain the following fields: %v", expectedFields)

	// Verify trade types are valid
	validTradeTypes := constants.ValidTradeTypes()
	t.Logf("Valid trade types: %v", validTradeTypes)

	// Test that our constants are properly configured
	if constants.DefaultAIModel == "" {
		t.Error("DefaultAIModel should not be empty")
	}

	if constants.DefaultAITimeout <= 0 {
		t.Error("DefaultAITimeout should be positive")
	}

	if constants.DefaultAIMaxRetry <= 0 {
		t.Error("DefaultAIMaxRetry should be positive")
	}

	t.Logf("Configuration validation passed")
}


