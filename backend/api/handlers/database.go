package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/internal/database"
)

// DatabaseHealthResponse represents the database health check response
type DatabaseHealthResponse struct {
	Status      string             `json:"status"`
	Timestamp   time.Time          `json:"timestamp"`
	Database    DatabaseStatusInfo `json:"database"`
	Performance PerformanceInfo    `json:"performance,omitempty"`
	Error       string             `json:"error,omitempty"`
}

// DatabaseStatusInfo contains database connection information
type DatabaseStatusInfo struct {
	Connected       bool `json:"connected"`
	MaxConnections  int  `json:"max_connections"`
	OpenConnections int  `json:"open_connections"`
	InUse           int  `json:"in_use"`
	Idle            int  `json:"idle"`
}

// PerformanceInfo contains database performance metrics
type PerformanceInfo struct {
	ResponseTimeMs int64 `json:"response_time_ms"`
}

// DatabaseHealthHandler handles database health check requests
func DatabaseHealthHandler(c *gin.Context) {
	startTime := time.Now()

	response := DatabaseHealthResponse{
		Timestamp: startTime,
		Status:    "healthy",
	}

	// Check if database is available
	if database.DB == nil {
		response.Status = "unhealthy"
		response.Error = "database connection is nil"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Perform health check
	sqlDB, err := database.DB.DB()
	if err != nil {
		response.Status = "unhealthy"
		response.Error = "failed to get underlying sql.DB: " + err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Test database connectivity with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		response.Status = "unhealthy"
		response.Error = "database ping failed: " + err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Get connection stats
	stats := sqlDB.Stats()
	response.Database = DatabaseStatusInfo{
		Connected:       true,
		MaxConnections:  stats.MaxOpenConnections,
		OpenConnections: stats.OpenConnections,
		InUse:           stats.InUse,
		Idle:            stats.Idle,
	}

	// Calculate response time
	responseTime := time.Since(startTime)
	response.Performance = PerformanceInfo{
		ResponseTimeMs: responseTime.Milliseconds(),
	}

	c.JSON(http.StatusOK, response)
}
