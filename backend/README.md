# Core API

# API Endpoints

### Portfolio Endpoints

- `GET /api/v1/portfolio/summary` - Portfolio overview
- `GET /api/v1/portfolio/holdings` - All current holdings
- `GET /api/v1/portfolio/holdings/{symbol}` - Single holding details

### Transaction Endpoints

- `GET /api/v1/transactions` - Transaction history
- `POST /api/v1/transactions` - Create new transaction
- `PUT /api/v1/transactions/{id}` - Update transaction
- `DELETE /api/v1/transactions/{id}` - Delete transaction

### Price Service Endpoints

- `GET /api/v1/price/current` - Current stock prices
- `GET /api/v1/price/historical` - Historical price data

A RESTful Golang API for tracking financial transactions with AI-powered transaction extraction from images.

## Features

- RESTful API with Gin framework
- JWT Authentication with rate limiting
- AI-powered transaction extraction from screenshots
- Gemini AI integration for image processing
- OpenAPI documentation with Postman collection
- Comprehensive testing framework
- Centralized prompt management system

## Prerequisites

- Go 1.21 or higher
- Valid Gemini API key for AI features

## Getting Started

### Clone the repository

```bash
git clone https://github.com/your-username/transaction-tracker.git
cd transaction-tracker/backend
```

### Install dependencies

```bash
go mod download
```

### Set environment variables

```bash
export SERVER_ADDR=":8080"
export JWT_SECRET="your-secret-key"
export GEMINI_API_KEY="your-gemini-api-key"
export AI_MODEL="gemini-2.0-flash"
export AI_TIMEOUT="30"
export AI_MAX_RETRY="3"
```

**Required Environment Variables:**

- `GEMINI_API_KEY`: Your Google AI API key for Gemini access
- `JWT_SECRET`: Secret key for JWT token signing

**Optional Environment Variables:**

- `SERVER_ADDR`: Server binding address (default: `:8080`)
- `AI_MODEL`: Gemini model to use (default: `gemini-2.0-flash`)
- `AI_TIMEOUT`: AI request timeout in seconds (default: `30`)
- `AI_MAX_RETRY`: Maximum retry attempts for AI requests (default: `3`)

### Run the application

```bash
go run main.go
```

The server will start on http://localhost:8080.

## AI Features

### Transaction Extraction

The API can process transaction screenshots from various brokers and extract structured data including:

- Stock ticker symbols and company names
- Trade dates and types (Buy/Sell/Dividends)
- Quantities and prices
- Exchange and currency information

### Supported Image Formats

- PNG, JPEG, GIF, WebP
- Multiple images per request
- Automatic image format detection

### AI Configuration

The system uses Google's Gemini AI model with configurable parameters:

- **Model**: Gemini 2.0 Flash (fast, cost-effective)
- **Timeout**: Configurable request timeout
- **Retry Logic**: Automatic retry with exponential backoff
- **Prompt Management**: Centralized, version-controlled prompts

## Testing

### Run All Tests

```bash
go test ./...
```

### Run AI Module Tests

```bash
go test ./test/ai/...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Integration Testing

For full integration tests with real AI API calls:

```bash
export GEMINI_API_KEY="your-real-api-key"
go test -v ./test/ai/ -tags=integration
```

## Project Structure

```
backend/
├── api/
│   ├── handlers/         # HTTP request handlers
│   ├── middlewares/      # Authentication & rate limiting
│   └── routes/           # Route definitions
├── config/               # Configuration management
├── docs/                 # API documentation
│   ├── api.yaml          # OpenAPI specification
│   └── postman_collection.json
├── internal/
│   ├── ai/               # AI client implementations
│   ├── prompts/          # AI prompt templates
│   └── utils/            # Utility functions
└── test/                 # Test files and test data
    ├── ai/               # AI module tests
    └── dummy-data/       # Test images and data
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Default**: 100 requests per minute per client
- **Configurable**: Adjust limits in environment variables
- **Headers**: Rate limit info included in response headers

## Authentication

JWT-based authentication with the following flow:

1. **Login**: POST to `/api/v1/login` with credentials
2. **Token**: Receive JWT token in response
3. **Access**: Include token in Authorization header: `Authorization: Bearer <token>`
4. **Expiry**: Tokens expire after 24 hours (configurable)
