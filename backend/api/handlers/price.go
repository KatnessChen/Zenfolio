package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/provider"
)

type PriceHandler struct {
	priceServiceManager *provider.PriceServiceManager
}

func NewPriceHandler(cfg *config.Config) *PriceHandler {
	manager := provider.NewPriceServiceManager(cfg)
	return &PriceHandler{
		priceServiceManager: manager,
	}
}

// GetCurrentPrice godoc
// @Summary Get current price for a single symbol
// @Description Get the current price for a specific stock symbol
// @Tags prices
// @Accept json
// @Produce json
// @Param symbol path string true "Stock symbol (e.g., AAPL)"
// @Success 200 {object} provider.PriceResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/prices/{symbol} [get]
func (h *PriceHandler) GetCurrentPrice(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Symbol parameter is required",
		})
		return
	}

	// Convert to uppercase for consistency
	symbol = strings.ToUpper(symbol)

	price, err := h.priceServiceManager.GetCurrentPrice(c.Request.Context(), symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch price data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    price,
	})
}

// GetCurrentPrices godoc
// @Summary Get current prices for multiple symbols
// @Description Get current prices for multiple stock symbols
// @Tags prices
// @Accept json
// @Produce json
// @Param symbols query string true "Comma-separated list of stock symbols (e.g., AAPL,GOOGL,MSFT)"
// @Success 200 {object} provider.BatchPriceResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/prices [get]
func (h *PriceHandler) GetCurrentPrices(c *gin.Context) {
	symbolsParam := c.Query("symbols")
	if symbolsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "symbols query parameter is required",
		})
		return
	}

	// Split symbols and clean them
	symbols := strings.Split(symbolsParam, ",")
	cleanSymbols := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		cleaned := strings.TrimSpace(strings.ToUpper(symbol))
		if cleaned != "" {
			cleanSymbols = append(cleanSymbols, cleaned)
		}
	}

	if len(cleanSymbols) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one valid symbol is required",
		})
		return
	}

	prices, err := h.priceServiceManager.GetCurrentPrices(c.Request.Context(), cleanSymbols)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch price data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prices,
	})
}

// GetPriceServiceHealth godoc
// @Summary Check price service health
// @Description Check if the price service is healthy and accessible
// @Tags prices
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/prices/health [get]
func (h *PriceHandler) GetPriceServiceHealth(c *gin.Context) {
	health, err := h.priceServiceManager.HealthCheck(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Price service is not healthy",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    health,
	})
}
