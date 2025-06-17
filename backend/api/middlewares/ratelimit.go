package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/config"
	"golang.org/x/time/rate"
)

// ClientRateLimiter stores rate limiters for different clients
type ClientRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	config   *config.Config
}

// NewClientRateLimiter creates a new ClientRateLimiter
func NewClientRateLimiter(cfg *config.Config) *ClientRateLimiter {
	return &ClientRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   cfg,
	}
}

// GetLimiter returns the rate limiter for a particular client
func (c *ClientRateLimiter) GetLimiter(clientID string) *rate.Limiter {
	c.mu.Lock()
	defer c.mu.Unlock()

	limiter, exists := c.limiters[clientID]
	if !exists {
		// Create a new rate limiter for the client
		limiter = rate.NewLimiter(rate.Limit(float64(c.config.RateLimitRequests)/c.config.RateLimitDuration.Seconds()), c.config.RateLimitRequests)
		c.limiters[clientID] = limiter
	}

	return limiter
}

// RateLimitMiddleware returns a middleware for rate limiting
func RateLimitMiddleware(rateLimiter *ClientRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as client identifier
		clientID := c.ClientIP()

		// For authenticated routes, you can use the user ID as the client identifier
		if userID, exists := c.Get("userID"); exists {
			if id, ok := userID.(string); ok {
				clientID = id
			}
		}

		// Get the rate limiter for this client
		limiter := rateLimiter.GetLimiter(clientID)

		// Check if the request is allowed
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
