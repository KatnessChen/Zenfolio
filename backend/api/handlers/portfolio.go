package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/transaction-tracker/backend/internal/models"
	"github.com/transaction-tracker/backend/internal/services"
)

// PortfolioHandler handles portfolio-related HTTP requests
type PortfolioHandler struct {
	portfolioService *services.PortfolioService
}

// NewPortfolioHandler creates a new portfolio handler
func NewPortfolioHandler(portfolioService *services.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioService: portfolioService,
	}
}

// GetSingleHoldingBasicInfo handles GET /api/v1/portfolio/holdings/{symbol}
func (h *PortfolioHandler) GetSingleHoldingBasicInfo(c *gin.Context) {
	// Get user ID from JWT token (from auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Get symbol from URL parameter
	symbol := strings.TrimSpace(strings.ToUpper(c.Param("symbol")))
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Symbol parameter is required",
		})
		return
	}

	// Validate symbol format (basic validation)
	if len(symbol) < 1 || len(symbol) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Symbol must be 1-10 characters",
		})
		return
	}

	// Get analysis type (default to basic)
	analysisType := c.DefaultQuery("analysis_type", "basic")
	if analysisType != "basic" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Only 'basic' analysis type is currently supported",
		})
		return
	}

	// Get stock basic info from service
	holdingInfo, err := h.portfolioService.GetSingleHoldingBasicInfo(c.Request.Context(), userID, symbol)
	if err != nil {
		if strings.Contains(err.Error(), "no transactions found") || strings.Contains(err.Error(), "no current holdings") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "failed to get current price") {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "Unable to fetch current price data",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get stock information",
		})
		return
	}

	// Create response
	response := models.SingleHoldingResponse{
		SingleHolding: *holdingInfo,
		Timestamp:     time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Stock basic info retrieved successfully",
		"data":    response,
	})
}

// GetAllHoldings handles GET /api/v1/portfolio/holdings
func (h *PortfolioHandler) GetAllHoldings(c *gin.Context) {
	// Get user ID from JWT token (from auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User ID not found in token",
		})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Get all holdings from service
	holdings, err := h.portfolioService.GetAllHoldings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get holdings information",
		})
		return
	}

	// Create response
	response := models.AllHoldingsResponse{
		Holdings:  holdings,
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}
