package constants

// AI Model Configuration
const (
	DefaultAIModel = "gemini-2.0-flash"
	DefaultAITimeout = 30
	DefaultAIMaxRetry = 3
)

// API Routes and Endpoints
const (
	APIVersion = "/api/v1"
	HealthEndpoint = "/health"
	LoginEndpoint = "/login"
	HelloWorldEndpoint = "/hello-world"
	ExtractTransEndpoint = "/extract-transactions"
)

// HTTP Headers
const (
	AuthorizationHeader = "Authorization"
	ContentTypeHeader = "Content-Type"
	BearerTokenPrefix = "Bearer"
)

// MIME Types
const (
	MimeTypePNG = "image/png"
	MimeTypeJPEG = "image/jpeg"
	MimeTypeGIF = "image/gif"
	MimeTypeWebP = "image/webp"
	
	MimeTypeJSON = "application/json"
	MimeTypeForm = "multipart/form-data"
)

// Transaction Types
const (
	TradeTypeBuy = "Buy"
	TradeTypeSell = "Sell"
	TradeTypeDividends = "Dividends"
)

// Common Exchanges
const (
	ExchangeNASDAQ = "NASDAQ"
	ExchangeNYSE = "NYSE"
	ExchangeOTC = "OTC"
)

// Common Currencies
const (
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	CurrencyGBP = "GBP"
	CurrencyJPY = "JPY"
	CurrencyCAD = "CAD"
	CurrencyTWD = "TWD"
)

// Error Messages
const (
	ErrMsgAuthHeaderRequired = "Authorization header is required"
	ErrMsgInvalidAuthFormat = "Invalid Authorization header format"
	ErrMsgInvalidToken = "Invalid token"
	ErrMsgTokenExpired = "Token expired"
	ErrMsgInvalidSigningMethod = "Invalid signing method"
	
	ErrMsgNoImagesProvided = "No images provided"
	ErrMsgImageProcessingFailed = "Image processing failed"
	ErrMsgAIRequestFailed = "AI request failed"
	
	ErrMsgInvalidTradeType = "Invalid trade type, must be Buy, Sell, or Dividends"
	ErrMsgTickerRequired = "Ticker should not be empty"
	ErrMsgTradeDateRequired = "TradeDate should not be empty"
	ErrMsgTradeTypeRequired = "TradeType should not be empty"
	ErrMsgNegativePrice = "Price should not be negative"
	
	ErrMsgInternalServer = "Internal server error"
	ErrMsgBadRequest = "Bad request"
	ErrMsgUnauthorized = "Unauthorized"
	ErrMsgForbidden = "Forbidden"
	ErrMsgNotFound = "Not found"
	ErrMsgTooManyRequests = "Too many requests"
)

// Success Messages
const (
	MsgTransactionsExtracted = "Transactions extracted successfully"
	MsgHealthCheckOK = "API is healthy"
	MsgLoginSuccessful = "Login successful"
	MsgHelloWorld = "Hello, World! You are authenticated."
)

// Environment Variable Names
const (
	EnvServerAddr = "SERVER_ADDR"
	EnvJWTSecret = "JWT_SECRET"
	EnvGeminiAPIKey = "GEMINI_API_KEY"
	EnvAIModel = "AI_MODEL"
	EnvAITimeout = "AI_TIMEOUT"
	EnvAIMaxRetry = "AI_MAX_RETRY"
)

// Default Values
const (
	DefaultServerAddr = ":8080"
	DefaultJWTExpiry = 24
)

// Rate Limiting
const (
	DefaultRateLimit = 100
	DefaultRateLimitWindow = 60
)

// File Upload Limits
const (
	MaxFileSize = 10 << 20
	MaxFilesPerBatch = 10
)

// Validation Constants
const (
	MinPasswordLength = 6
	MaxUsernameLength = 50
	MaxTickerLength = 10
)

// ValidTradeTypes returns a slice of valid trade types
func ValidTradeTypes() []string {
	return []string{TradeTypeBuy, TradeTypeSell, TradeTypeDividends}
}

// ValidTradeTypesMap returns a map of valid trade types for quick lookup
func ValidTradeTypesMap() map[string]bool {
	return map[string]bool{
		TradeTypeBuy:       true,
		TradeTypeSell:      true,
		TradeTypeDividends: true,
	}
}

// SupportedImageMimeTypes returns a slice of supported image MIME types
func SupportedImageMimeTypes() []string {
	return []string{MimeTypePNG, MimeTypeJPEG, MimeTypeGIF, MimeTypeWebP}
}

// SupportedImageMimeTypesMap returns a map of supported image MIME types for quick lookup
func SupportedImageMimeTypesMap() map[string]bool {
	return map[string]bool{
		MimeTypePNG:  true,
		MimeTypeJPEG: true,
		MimeTypeGIF:  true,
		MimeTypeWebP: true,
	}
}
