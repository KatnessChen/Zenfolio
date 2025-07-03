# Testing Framework

This directory contains comprehensive tests for the Transaction Tracker backend, with a focus on AI module functionality and integration testing.

## Test Structure

```
test/
├── ai/                    # AI module tests
│   └── client_test.go    # AI client interface and implementations
├── dummy-data/           # Test data and sample files
│   └── transaction-screenshots/  # Sample transaction images
│       └── Firstrade_1.png
└── README.md            # This file
```

## Running Tests

### All Tests

```bash
# Run all tests in the project
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage report
go test -cover ./...
```

### AI Module Tests

```bash
# Run all AI tests
go test ./test/ai/...

# Run with verbose output
go test -v ./test/ai/...

# Run specific test function
go test ./test/ai -run TestGeminiClientCreation

# Run with race detection
go test -race ./test/ai/...
```

### Integration Tests

```bash
# Run integration tests (requires real API key)
export GEMINI_API_KEY="your-actual-api-key"
go test -v ./test/ai/ -tags=integration

# Run integration tests with timeout
go test -v ./test/ai/ -tags=integration -timeout=60s
```

## Test Categories

### Unit Tests

Test individual components in isolation:

- **AI Client Interface**: Tests for the `ai.Client` interface
- **Factory Pattern**: Tests for AI client factory functions
- **Configuration**: Tests for AI configuration validation
- **Prompt Loading**: Tests for prompt template loading

### Integration Tests

Test full AI workflows:

- **Image Processing**: End-to-end transaction extraction from real images
- **API Integration**: Full HTTP request/response cycle testing
- **Error Handling**: Network failures, invalid API keys, malformed responses

### Mock Tests

Tests using mocked dependencies:

- **Network Failures**: Simulated API downtime
- **Rate Limiting**: API rate limit scenarios
- **Invalid Responses**: Malformed AI responses

## Test Data

### Sample Images

The `dummy-data/transaction-screenshots/` directory contains:

- **Firstrade_1.png**: Sample Firstrade brokerage screenshot
- **Future additions**: Screenshots from other brokers

These images are used for:

- Integration testing with real AI models
- Regression testing of extraction accuracy
- Performance benchmarking

### Test Configuration

Tests use environment variables for configuration:

```bash
# Required for integration tests
export GEMINI_API_KEY="your-api-key"

# Optional test configuration
export TEST_AI_MODEL="gemini-2.0-flash"
export TEST_TIMEOUT="30"
export TEST_MAX_RETRY="3"
```

## Environment Variables

### Required for Integration Tests

```bash
export GEMINI_API_KEY="your-actual-api-key"
```

### Optional Test Configuration

```bash
export TEST_AI_MODEL="gemini-2.0-flash"      # AI model to use for tests
export TEST_TIMEOUT="30"                      # Request timeout in seconds
export TEST_MAX_RETRY="3"                     # Maximum retry attempts
export TEST_VERBOSE="true"                    # Enable verbose test output
```

## Test Best Practices

### Writing New Tests

1. **Use Table-Driven Tests**: For testing multiple scenarios
2. **Mock External Dependencies**: Use interfaces and dependency injection
3. **Test Error Conditions**: Include both success and failure cases
4. **Use Descriptive Names**: Test function names should describe what they test
5. **Clean Up Resources**: Ensure tests don't leave artifacts

### Example Test Structure

```go
func TestAIClientExtraction(t *testing.T) {
    tests := []struct {
        name           string
        inputImages    []ai.FileInput
        expectedResult *ai.ExtractResponse
        expectError    bool
    }{
        {
            name: "successful extraction",
            inputImages: []ai.FileInput{
                {Data: mockImageData, Filename: "test.png", MimeType: "image/png"},
            },
            expectError: false,
        },
        {
            name: "invalid image data",
            inputImages: []ai.FileInput{
                {Data: invalidData, Filename: "bad.png", MimeType: "image/png"},
            },
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Continuous Integration

Tests are designed to run in CI/CD environments:

- **Unit Tests**: Run without external dependencies
- **Integration Tests**: Run with valid API keys in secure environments
- **Coverage Reports**: Generate coverage reports for monitoring
- **Performance Tests**: Track AI response times and accuracy

## Troubleshooting

### Common Issues

1. **API Key Errors**: Ensure `GEMINI_API_KEY` is set for integration tests
2. **Network Timeouts**: Increase timeout values for slow connections
3. **Rate Limiting**: Space out test runs to avoid API rate limits
4. **File Not Found**: Ensure test data files are in the correct location

### Debug Mode

Enable debug logging for tests:

```bash
export DEBUG=true
go test -v ./test/ai/
```

Note: Unit tests don't require real API keys and will work with dummy values. Integration tests require a valid Gemini API key for full functionality testing.
