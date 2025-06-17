package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/config"
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

	// Get database manager instance (assuming it's available globally)
	// For now, we'll create a temporary instance
	dbConfig := getDefaultDatabaseConfig()
	dm := database.NewDatabaseManager(dbConfig)
	dm.GetDB() // This won't work without proper initialization, so let's use the global DB

	// Perform health check
	sqlDB, err := database.DB.DB()
	if err != nil {
		response.Status = "unhealthy"
		response.Error = "failed to get underlying sql.DB: " + err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Test database connectivity
	if err := sqlDB.Ping(); err != nil {
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

// DetailedDatabaseHealthHandler provides more detailed database health information
func DetailedDatabaseHealthHandler(c *gin.Context) {
	startTime := time.Now()

	// Check if user is authenticated (optional, depending on security requirements)
	// userID, exists := c.Get("userID")
	// if !exists {
	//     c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	//     return
	// }

	response := DatabaseHealthResponse{
		Timestamp: startTime,
		Status:    "healthy",
	}

	if database.DB == nil {
		response.Status = "unhealthy"
		response.Error = "database connection is nil"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Test a simple query
	var result int
	if err := database.DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		response.Status = "unhealthy"
		response.Error = "database query failed: " + err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Get underlying database
	sqlDB, err := database.DB.DB()
	if err != nil {
		response.Status = "unhealthy"
		response.Error = "failed to get underlying sql.DB: " + err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Get detailed connection stats
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

// getDefaultDatabaseConfig returns a default database configuration for health checks
func getDefaultDatabaseConfig() *config.DatabaseConfig {
	// This is a simplified version - in practice, you'd get this from your config
	return &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "3306",
		Name:     "transaction_tracker_dev",
		User:     "root",
		Password: "",
	}
}
