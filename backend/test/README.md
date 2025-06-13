# AI Module Tests

This directory contains tests for the AI module components.

## Running Tests

To run all AI tests:

```bash
go test ./test/ai/...
```

To run tests with verbose output:

```bash
go test -v ./test/ai/...
```

To run specific test:

```bash
go test ./test/ai -run TestGeminiClientCreation
```

## Test Structure

- `client_test.go` - Tests for AI client interface and Gemini implementation

## Test Requirements

Tests require:

- Valid Go environment
- Gemini SDK dependencies (installed via `go mod download`)
- Tests use dummy API keys and expect authentication errors

## Environment Variables for Integration Tests

For integration tests with real API calls, set:

```bash
export GEMINI_API_KEY="your-actual-api-key"
```

Note: Unit tests don't require real API keys and will work with dummy values.
