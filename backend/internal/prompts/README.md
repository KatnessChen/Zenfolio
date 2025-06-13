# AI Prompts

This directory contains reusable prompt templates for AI model interactions.

## Files

### `system_instruction.txt`

The system instruction that defines the AI assistant's role and behavior. This is used to configure the AI model's general approach to processing requests.

### `transaction_extraction.txt`

The main prompt template used for extracting transaction data from screenshots. Includes:

- JSON output format specification
- Business rules for data extraction
- Data validation requirements

## Usage

The prompts are loaded using the `prompts` package:

```go
import "github.com/transaction-tracker/backend/internal/prompts"

// Load system instruction
systemInstruction, err := prompts.GetSystemInstruction()

// Load transaction extraction prompt
extractionPrompt, err := prompts.GetTransactionExtractionPrompt()
```

## Benefits of Centralized Prompts

1. **Maintainability**: Easy to update prompts without modifying Go code
2. **Consistency**: Same prompts used across different AI model implementations
3. **Version Control**: Prompt changes are tracked in Git
4. **Testing**: Prompts can be tested independently
5. **Reusability**: Prompts can be shared between different components

## Embedded Files

The prompt files are embedded into the Go binary using `go:embed`, so they don't need to be distributed separately with the application.
