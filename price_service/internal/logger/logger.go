package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var globalLogger zerolog.Logger

// H is a convenient type alias for map[string]interface{}
type H map[string]interface{}

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *strings.Builder
}

// InitLogger initializes the global logger for price-service with advanced configuration
func InitLogger() {
	// Read log level from environment variable, default to debug
	logLevel := zerolog.DebugLevel
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		if level, err := zerolog.ParseLevel(envLevel); err == nil {
			logLevel = level
		}
	}

	// Use standard output (no ConsoleWriter)
	var output io.Writer = os.Stdout

	// Configure global logger
	globalLogger = zerolog.New(output).
		Level(logLevel).
		With().
		Timestamp().
		Str("service", "price-service").
		Logger()

	// Set as global zerolog logger
	log.Logger = globalLogger

	// Log initialization
	globalLogger.Info().
		Str("level", logLevel.String()).
		Str("format", "json").
		Msg("Logger initialized successfully")
}

// extractHandlerName extracts the handler function name from Gin's handler name
func extractHandlerName(handlerName string) string {
	if handlerName == "" {
		return "unknown"
	}

	// Extract function name from full path
	parts := strings.Split(handlerName, ".")
	if len(parts) > 0 {
		// Get the last part which should be the function name
		funcName := parts[len(parts)-1]
		// Remove any suffix like "-fm" that gin adds
		if idx := strings.Index(funcName, "-"); idx != -1 {
			funcName = funcName[:idx]
		}
		return funcName
	}

	return "handler"
}

// extractCaller generates the caller field from request method and path
func extractCaller(c *gin.Context) string {
	return fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
}

// sanitizeData removes sensitive fields from data
func sanitizeData(data interface{}) interface{} {
	sensitiveFields := []string{"password", "token", "secret", "key", "auth", "credential"}

	switch v := data.(type) {
	case map[string]interface{}:
		sanitized := make(map[string]interface{})
		for key, value := range v {
			lowerKey := strings.ToLower(key)
			isSensitive := false
			for _, sensitive := range sensitiveFields {
				if strings.Contains(lowerKey, sensitive) {
					isSensitive = true
					break
				}
			}
			if isSensitive {
				sanitized[key] = "[REDACTED]"
			} else {
				sanitized[key] = sanitizeData(value)
			}
		}
		return sanitized
	case []interface{}:
		sanitized := make([]interface{}, len(v))
		for i, item := range v {
			sanitized[i] = sanitizeData(item)
		}
		return sanitized
	default:
		return data
	}
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// GinLoggerMiddleware returns a configured gin-contrib/logger middleware with advanced features
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Create a custom response writer to capture the response body
		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &strings.Builder{},
		}

		// Replace the response writer
		c.Writer = rw

		// Apply the gin-contrib/logger middleware
		loggerMiddleware := logger.SetLogger(
			logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
				return l.With().
					Str("service", "price-service").
					Str("scope", extractHandlerName(c.HandlerName())).
					Str("caller", extractCaller(c)).
					Time("timestamp", time.Now()).
					Logger()
			}),
			logger.WithUTC(true),
			logger.WithContext(func(c *gin.Context, e *zerolog.Event) *zerolog.Event {
				// Extract input data
				input := H{}

				// Add query parameters
				if len(c.Request.URL.RawQuery) > 0 {
					for key, values := range c.Request.URL.Query() {
						if len(values) > 0 {
							input["query_"+key] = sanitizeData(values[0])
						}
					}
				}

				// Add path parameters
				for _, param := range c.Params {
					input["param_"+param.Key] = sanitizeData(param.Value)
				}

				// Add request body for non-GET requests
				if c.Request.Method != "GET" && c.Request.ContentLength > 0 {
					if bodyBytes, err := io.ReadAll(c.Request.Body); err == nil {
						// Reset body for subsequent reading
						c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

						var bodyData interface{}
						if err := json.Unmarshal(bodyBytes, &bodyData); err == nil {
							input["body"] = sanitizeData(bodyData)
						} else {
							input["body"] = sanitizeData(string(bodyBytes))
						}
					}
				}

				// Capture response data
				output := H{}
				responseBody := rw.body.String()
				if responseBody != "" {
					var responseData interface{}
					if err := json.Unmarshal([]byte(responseBody), &responseData); err == nil {
						output["body"] = sanitizeData(responseData)
					} else {
						output["body"] = sanitizeData(responseBody)
					}
				}

				return e.Interface("input", input).
					Interface("output", output).
					Str("user_agent", c.Request.UserAgent()).
					Str("client_ip", c.ClientIP())
			}),
		)

		// Execute the logger middleware
		loggerMiddleware(c)
	}
}

// getFields extracts the first H map from variadic arguments or returns empty H
func getFields(fields ...H) H {
	if len(fields) > 0 {
		return fields[0]
	}
	return H{}
}

// Debug logs a debug level message
func Debug(message string, fields ...H) {
	globalLogger.Debug().Fields(map[string]interface{}(getFields(fields...))).Msg(message)
}

// Info logs an info level message
func Info(message string, fields ...H) {
	globalLogger.Info().Fields(map[string]interface{}(getFields(fields...))).Msg(message)
}

// Warn logs a warning level message
func Warn(message string, fields ...H) {
	globalLogger.Warn().Fields(map[string]interface{}(getFields(fields...))).Msg(message)
}

// Error logs an error level message
func Error(message string, err error, fields ...H) {
	logEvent := globalLogger.Error()
	if err != nil {
		logEvent = logEvent.Err(err)
	}
	logEvent.Fields(map[string]interface{}(getFields(fields...))).Msg(message)
}

// GetLogger returns the global logger instance
func GetLogger() *zerolog.Logger {
	return &globalLogger
}
