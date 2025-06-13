package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/transaction-tracker/backend/internal/prompts"
	"google.golang.org/api/option"
)

// GeminiClient implements the Client interface using Google's Gemini API
type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
	config *Config
}

// NewGeminiClient creates a new Gemini AI client
func NewGeminiClient(config *Config) (*GeminiClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required for Gemini client")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Set default model if not specified
	modelName := config.Model
	if modelName == "" {
		modelName = "gemini-2.5-flash"
	}

	model := client.GenerativeModel(modelName)
	
	// Configure model settings for better JSON responses
	model.ResponseMIMEType = "application/json"
	
	// Load system instruction from prompt file
	systemInstruction, err := prompts.LoadPrompt("system_instruction.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load system instruction: %w", err)
	}
	
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}

	return &GeminiClient{
		client: client,
		model:  model,
		config: config,
	}, nil
}

// ExtractTransactions processes images and extracts transaction data using Gemini
func (g *GeminiClient) ExtractTransactions(ctx context.Context, images []ImageInput) (*ExtractResponse, error) {
	if len(images) == 0 {
		return &ExtractResponse{
			Success: false,
			Message: "No images provided",
		}, nil
	}

	// Load the transaction extraction prompt
	prompt, promptErr := prompts.LoadPrompt("transaction_extraction.txt")
	if promptErr != nil {
		return &ExtractResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to load extraction prompt: %v", promptErr),
		}, fmt.Errorf("failed to load extraction prompt: %w", promptErr)
	}

	// Prepare the content parts for the request
	parts := []genai.Part{genai.Text(prompt)}

	// Add each image to the request
	for _, img := range images {
		imageData, err := io.ReadAll(img.Data)
		if err != nil {
			log.Printf("Failed to read image %s: %v", img.Filename, err)
			continue
		}

		parts = append(parts, genai.ImageData(img.MimeType, imageData))
	}

	// Set timeout if specified
	if g.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(g.config.Timeout)*time.Second)
		defer cancel()
	}

	// Generate content with retry logic
	var resp *genai.GenerateContentResponse
	var err error
	
	maxRetry := g.config.MaxRetry
	if maxRetry <= 0 {
		maxRetry = 3
	}

	for attempt := 0; attempt < maxRetry; attempt++ {
		resp, err = g.model.GenerateContent(ctx, parts...)
		if err == nil {
			break
		}
		
		if attempt < maxRetry-1 {
			log.Printf("Attempt %d failed, retrying: %v", attempt+1, err)
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}

	if err != nil {
		return &ExtractResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to generate content after %d attempts: %v", maxRetry, err),
		}, err
	}

	// Parse the response
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return &ExtractResponse{
			Success: false,
			Message: "No response received from AI model",
		}, nil
	}

	// Extract the text response
	responseText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	
	// Parse JSON response
	var result struct {
		Transactions []TransactionData `json:"transactions"`
	}

	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		log.Printf("Failed to parse AI response as JSON: %v\nResponse: %s", err, responseText)
		return &ExtractResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to parse AI response: %v", err),
		}, nil
	}

	return &ExtractResponse{
		Transactions: result.Transactions,
		Success:      true,
		Message:      fmt.Sprintf("Successfully extracted %d transactions", len(result.Transactions)),
	}, nil
}

// Health checks if the Gemini client is working properly
func (g *GeminiClient) Health(ctx context.Context) error {
	// Simple health check by making a minimal request
	parts := []genai.Part{genai.Text("Health check - please respond with 'OK'")}
	
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := g.model.GenerateContent(ctx, parts...)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	return nil
}

// Close closes the client and cleans up resources
func (g *GeminiClient) Close() error {
	if g.client != nil {
		return g.client.Close()
	}
	return nil
}
