package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/transaction-tracker/backend/internal/constants"
	"github.com/transaction-tracker/backend/internal/prompts"
	"google.golang.org/api/option"
)

// ModelType represents different AI model providers
type ModelType string

const (
	ModelTypeGemini ModelType = "gemini"
	// Future model types can be added here
	// ModelTypeOpenAI ModelType = "openai"
	// ModelTypeClaude ModelType = "claude"
)

// AIModelClient is a generic AI model client that can work with different providers
type AIModelClient struct {
	modelType ModelType
	config    *Config
	
	// Provider-specific clients (only one will be active at a time)
	geminiClient *genai.Client
	geminiModel  *genai.GenerativeModel
	
	// Future provider clients can be added here
	// openaiClient *openai.Client
	// claudeClient *claude.Client
}

// NewAIModelClient creates a new AI model client based on the configuration
func NewAIModelClient(config *Config) (*AIModelClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required for AI model client")
	}

	// Determine model type from config
	modelType := determineModelType(config.Model)
	
	client := &AIModelClient{
		modelType: modelType,
		config:    config,
	}

	// Initialize the appropriate provider client
	switch modelType {
	case ModelTypeGemini:
		if err := client.initGeminiClient(); err != nil {
			return nil, fmt.Errorf("failed to initialize Gemini client: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported model type: %s", modelType)
	}

	return client, nil
}

// determineModelType determines the AI provider based on the model name
func determineModelType(modelName string) ModelType {
	if modelName == "" {
		modelName = constants.DefaultAIModel
	}

	// Check for Gemini models
	switch modelName {
	case "gemini-2.0-flash", "gemini-1.5-pro", "gemini-1.5-flash":
		return ModelTypeGemini
	}

	// Future model detection logic can be added here
	// if strings.HasPrefix(modelName, "gpt-") {
	//     return ModelTypeOpenAI
	// }
	// if strings.HasPrefix(modelName, "claude-") {
	//     return ModelTypeClaude
	// }

	// Default to Gemini for backward compatibility
	return ModelTypeGemini
}

// initGeminiClient initializes the Gemini client
func (c *AIModelClient) initGeminiClient() error {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(c.config.APIKey))
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Set default model if not specified
	modelName := c.config.Model
	if modelName == "" {
		modelName = constants.DefaultAIModel
	}

	model := client.GenerativeModel(modelName)
	
	// Configure model settings for better JSON responses
	model.ResponseMIMEType = constants.MimeTypeJSON
	
	// Load system instruction from prompt file
	systemInstruction, err := prompts.LoadPrompt("system_instruction.txt")
	if err != nil {
		return fmt.Errorf("failed to load system instruction: %w", err)
	}
	
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}

	c.geminiClient = client
	c.geminiModel = model
	
	return nil
}

// ExtractTransactions processes images and extracts transaction data using the configured AI model
func (c *AIModelClient) ExtractTransactions(ctx context.Context, images []ImageInput) (*ExtractResponse, error) {
	if len(images) == 0 {
		return &ExtractResponse{
			Success: false,
			Message: constants.ErrMsgNoImagesProvided,
		}, nil
	}

	// Route to the appropriate implementation based on model type
	switch c.modelType {
	case ModelTypeGemini:
		return c.extractTransactionsGemini(ctx, images)
	default:
		return &ExtractResponse{
			Success: false,
			Message: fmt.Sprintf("Unsupported model type: %s", c.modelType),
		}, fmt.Errorf("unsupported model type: %s", c.modelType)
	}
}

// extractTransactionsGemini handles transaction extraction using Gemini models
func (c *AIModelClient) extractTransactionsGemini(ctx context.Context, images []ImageInput) (*ExtractResponse, error) {
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
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(c.config.Timeout)*time.Second)
		defer cancel()
	}

	// Generate content with retry logic
	var resp *genai.GenerateContentResponse
	var err error
	
	maxRetry := c.config.MaxRetry
	if maxRetry <= 0 {
		maxRetry = constants.DefaultAIMaxRetry
	}

	for attempt := 0; attempt < maxRetry; attempt++ {
		resp, err = c.geminiModel.GenerateContent(ctx, parts...)
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
	
	return c.parseTransactionResponse(responseText)
}

// parseTransactionResponse parses the AI response into transaction data (generic for all models)
func (c *AIModelClient) parseTransactionResponse(responseText string) (*ExtractResponse, error) {
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
		Message:      constants.MsgTransactionsExtracted,
	}, nil
}

// Health checks if the AI model client is working properly
func (c *AIModelClient) Health(ctx context.Context) error {
	switch c.modelType {
	case ModelTypeGemini:
		return c.healthCheckGemini(ctx)
	default:
		return fmt.Errorf("health check not implemented for model type: %s", c.modelType)
	}
}

// healthCheckGemini performs a health check for Gemini models
func (c *AIModelClient) healthCheckGemini(ctx context.Context) error {
	// Simple health check by making a minimal request
	parts := []genai.Part{genai.Text("Health check - please respond with 'OK'")}
	
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.geminiModel.GenerateContent(ctx, parts...)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	return nil
}

// Close closes the client and cleans up resources
func (c *AIModelClient) Close() error {
	switch c.modelType {
	case ModelTypeGemini:
		if c.geminiClient != nil {
			return c.geminiClient.Close()
		}
	}
	// Future model cleanup logic can be added here
	return nil
}

// GetModelType returns the current model type
func (c *AIModelClient) GetModelType() ModelType {
	return c.modelType
}

// GetModelName returns the configured model name
func (c *AIModelClient) GetModelName() string {
	return c.config.Model
}
