package ai_test

import (
	"context"
	"strings"
	"testing"

	"github.com/transaction-tracker/backend/internal/ai"
)

func TestGeminiClientCreation(t *testing.T) {
	// Test with empty API key
	config := &ai.Config{}
	_, err := ai.NewGeminiClient(config)
	if err == nil {
		t.Error("Expected error when API key is empty")
	}

	// Test with valid config (but dummy API key)
	config = &ai.Config{
		APIKey:   "dummy-key-for-testing",
		Model:    "gemini-2.5-flash", // As of mid-2025; check https://ai.google.dev/models for latest models
		Timeout:  30,
		MaxRetry: 3,
	}

	client, err := ai.NewGeminiClient(config)
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
		Model:    "gemini-2.5-flash",
		Timeout:  30,
		MaxRetry: 3,
	}

	client, err := ai.NewGeminiClient(config)
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

	if !strings.Contains(resp.Message, "No images provided") {
		t.Errorf("Expected 'No images provided' message, got: %s", resp.Message)
	}
}
