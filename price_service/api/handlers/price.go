package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/price_service/internal/cache"
	"github.com/transaction-tracker/price_service/internal/config"
	"github.com/transaction-tracker/price_service/internal/models"
	"github.com/transaction-tracker/price_service/internal/provider"
)

type PriceHandler struct {
	cache    *cache.Service
	provider provider.StockPriceProvider
	config   *config.Config
}

func NewPriceHandler(cache *cache.Service, provider provider.StockPriceProvider, config *config.Config) *PriceHandler {
	return &PriceHandler{
		cache:    cache,
		provider: provider,
		config:   config,
	}
}

// GetCurrentPrices handles GET /api/v1/price/current/symbols
func (h *PriceHandler) GetCurrentPrices(c *gin.Context) {
	symbolsParam := c.Query("symbols")
	if symbolsParam == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error: models.ErrorDetail{
				Code:    models.ErrInvalidInput,
				Message: "symbols parameter is required",
			},
		})
		return
	}

	symbols := strings.Split(symbolsParam, ",")
	if len(symbols) > h.config.Cache.MaxSymbolsPerReq {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error: models.ErrorDetail{
				Code:    models.ErrInvalidInput,
				Message: "too many symbols requested",
			},
		})
		return
	}

	// Clean and validate symbols
	var validSymbols []string
	for _, symbol := range symbols {
		cleanSymbol := strings.TrimSpace(strings.ToUpper(symbol))
		if cleanSymbol != "" {
			validSymbols = append(validSymbols, cleanSymbol)
		}
	}

	if len(validSymbols) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error: models.ErrorDetail{
				Code:    models.ErrInvalidInput,
				Message: "no valid symbols provided",
			},
		})
		return
	}

	var result []models.SymbolCurrentPrice
	var missingSymbols []string

	// Check cache first
	for _, symbol := range validSymbols {
		cached, err := h.cache.GetCurrentPrice(c.Request.Context(), symbol)
		if err != nil {
			// Log error but continue
			missingSymbols = append(missingSymbols, symbol)
		} else if cached != nil {
			result = append(result, *cached)
		} else {
			missingSymbols = append(missingSymbols, symbol)
		}
	}

	// Fetch missing symbols from provider
	if len(missingSymbols) > 0 {
		fetchedPrices, err := h.provider.GetCurrentPrices(c.Request.Context(), missingSymbols)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
				Success: false,
				Error: models.ErrorDetail{
					Code:    models.ErrServiceUnavailable,
					Message: "failed to fetch price data",
				},
			})
			return
		}

		// Cache the fetched prices and add to result
		for _, price := range fetchedPrices {
			if err := h.cache.SetCurrentPrice(c.Request.Context(), price.Symbol, &price); err != nil {
				// Log error but continue
				log.Printf("error caching current price for %s: %v", price.Symbol, err)
			}
			result = append(result, price)
		}
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success:   true,
		Data:      result,
		Timestamp: time.Now(),
	})
}

// GetHistoricalPrices handles GET /api/v1/price/historical/symbol
func (h *PriceHandler) GetHistoricalPrices(c *gin.Context) {
	symbol := strings.TrimSpace(strings.ToUpper(c.Query("symbol")))
	if symbol == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error: models.ErrorDetail{
				Code:    models.ErrInvalidInput,
				Message: "symbol parameter is required",
			},
		})
		return
	}

	// Parse date-related parameters
	dateParam := strings.TrimSpace(c.Query("date"))
	fromParam := strings.TrimSpace(c.Query("from"))
	toParam := strings.TrimSpace(c.Query("to"))

	// Validate parameter combinations
	if err := h.validateDateParameters(dateParam, fromParam, toParam); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error: models.ErrorDetail{
				Code:    models.ErrInvalidInput,
				Message: err.Error(),
			},
		})
		return
	}

	// Handle single date query
	if dateParam != "" {
		h.handleSingleDateQuery(c, symbol, dateParam)
		return
	}

	// Handle date range query
	if fromParam != "" && toParam != "" {
		h.handleDateRangeQuery(c, symbol, fromParam, toParam)
		return
	}

	// Handle resolution-based query
	h.handleResolutionQuery(c, symbol)
}

// handleSingleDateQuery processes single date historical price requests
func (h *PriceHandler) handleSingleDateQuery(c *gin.Context, symbol, dateParam string) {
	// Validate date format and future date using the shared validation method
	if err := h.validateSingleDate(dateParam); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error: models.ErrorDetail{
				Code:    models.ErrInvalidInput,
				Message: err.Error(),
			},
		})
		return
	}

	// Adjust the requested date to the last trading day
	requestedDate, _ := time.Parse("2006-01-02", dateParam)
	adjustedDate := h.getLastTradingDay(requestedDate)
	adjustedDateStr := adjustedDate.Format("2006-01-02")

	// Check if we have daily historical data cached
	cachedData, err := h.cache.GetHistoricalPrice(c.Request.Context(), symbol, models.ResolutionDaily)
	if err == nil && cachedData != nil && len(cachedData.HistoricalPrices) > 0 {
		// Get the most recent price record (first record since sorted newest to oldest)
		latestPriceRecord := cachedData.HistoricalPrices[0]
		latestDate, err := time.Parse("2006-01-02", latestPriceRecord.Date)
		if err == nil {
			// If the latest record is on or after the adjusted date, we have the data in cache
			if latestDate.Equal(adjustedDate) || latestDate.After(adjustedDate) {
				// Find the exact date in cache
				for _, priceData := range cachedData.HistoricalPrices {
					if priceData.Date == adjustedDateStr {
						result := &models.SymbolHistoricalPrice{
							Symbol:     symbol,
							Resolution: models.ResolutionDaily,
							HistoricalPrices: []models.ClosePrice{
								{
									Date:  adjustedDateStr,
									Price: priceData.Price,
								},
							},
						}

						c.JSON(http.StatusOK, models.SuccessResponse{
							Success: true,
							Data:    result,
						})
						return
					}
				}
			}
		}

		// If latest record is before the adjusted date, invalidate cache and fetch fresh data
		if err := h.cache.DeleteHistoricalPrice(c.Request.Context(), symbol, models.ResolutionDaily); err != nil {
			log.Printf("error invalidating cache for %s: %v", symbol, err)
		}
	}

	// Fetch from provider (either no cache or cache needs update)
	historicalData, err := h.provider.GetHistoricalPrices(c.Request.Context(), symbol, models.ResolutionDaily)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Success: false,
			Error: models.ErrorDetail{
				Code:    models.ErrServiceUnavailable,
				Message: "failed to fetch historical data for date",
			},
		})
		return
	}

	// Cache the fresh data
	if err := h.cache.SetHistoricalPrice(c.Request.Context(), symbol, models.ResolutionDaily, historicalData); err != nil {
		log.Printf("error caching historical price for %s: %v", symbol, err)
	}

	// Find the adjusted date in the fresh data
	for _, priceData := range historicalData.HistoricalPrices {
		if priceData.Date == adjustedDateStr {
			result := &models.SymbolHistoricalPrice{
				Symbol:     symbol,
				Resolution: models.ResolutionDaily,
				HistoricalPrices: []models.ClosePrice{
					{
						Date:  adjustedDateStr,
						Price: priceData.Price,
					},
				},
			}

			c.JSON(http.StatusOK, models.SuccessResponse{
				Success:   true,
				Data:      result,
				Timestamp: time.Now(),
			})
			return
		}
	}

	// If we still can't find the adjusted date, return error
	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Success: false,
		Error: models.ErrorDetail{
			Code:    models.ErrSymbolNotFound,
			Message: fmt.Sprintf("no data available for trading date %s", adjustedDateStr),
		},
	})
}

// handleDateRangeQuery processes date range historical price requests
func (h *PriceHandler) handleDateRangeQuery(c *gin.Context, symbol, fromParam, toParam string) {
	// Check if we have daily historical data cached that might contain our date range
	cachedData, err := h.cache.GetHistoricalPrice(c.Request.Context(), symbol, models.ResolutionDaily)
	if err == nil && cachedData != nil {
		// Check cache coverage for the requested date range
		coverage := h.checkCacheCoverage(cachedData, fromParam, toParam)

		switch coverage {
		case "full":
			// Cache covers full range → filter and return
			filteredData := h.filterHistoricalDataByDateParams(cachedData, "", fromParam, toParam)
			if len(filteredData.HistoricalPrices) > 0 {
				c.JSON(http.StatusOK, models.SuccessResponse{
					Success: true,
					Data:    filteredData,
				})
				return
			}
		case "partial":
			// Cache covers partial range → invalidate cache, fetch fresh data

			// Invalidate the existing cache to ensure we get fresh data from provider
			if err := h.cache.DeleteHistoricalPrice(c.Request.Context(), symbol, models.ResolutionDaily); err != nil {
				log.Printf("error invalidating cache for %s: %v", symbol, err)
			}

			freshData, err := h.provider.GetHistoricalPrices(c.Request.Context(), symbol, models.ResolutionDaily)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
					Success: false,
					Error: models.ErrorDetail{
						Code:    models.ErrServiceUnavailable,
						Message: "failed to fetch historical data for date range",
					},
				})
				return
			}

			// Update cache with fresh data
			if err := h.cache.SetHistoricalPrice(c.Request.Context(), symbol, models.ResolutionDaily, freshData); err != nil {
				log.Printf("error updating cache for %s: %v", symbol, err)
			}

			// Filter for the requested date range
			filteredData := h.filterHistoricalDataByDateParams(freshData, "", fromParam, toParam)
			c.JSON(http.StatusOK, models.SuccessResponse{
				Success: true,
				Data:    filteredData,
			})
			return
		}
	}

	// If no cache coverage → fetch full daily data and cache
	historicalData, err := h.provider.GetHistoricalPrices(c.Request.Context(), symbol, models.ResolutionDaily)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Success: false,
			Error: models.ErrorDetail{
				Code:    models.ErrServiceUnavailable,
				Message: "failed to fetch historical data for date range",
			},
		})
		return
	}

	// Cache the full daily data
	if err := h.cache.SetHistoricalPrice(c.Request.Context(), symbol, models.ResolutionDaily, historicalData); err != nil {
		log.Printf("error caching historical price for %s: %v", symbol, err)
	}

	// Filter for the requested date range
	filteredData := h.filterHistoricalDataByDateParams(historicalData, "", fromParam, toParam)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success:   true,
		Data:      filteredData,
		Timestamp: time.Now(),
	})
}

// handleResolutionQuery processes resolution-based historical price requests
func (h *PriceHandler) handleResolutionQuery(c *gin.Context, symbol string) {
	resolutionStr := c.DefaultQuery("resolution", "daily")

	// Validate resolution
	resolution := models.Resolution(resolutionStr)
	if resolution != models.ResolutionDaily &&
		resolution != models.ResolutionWeekly &&
		resolution != models.ResolutionMonthly {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error: models.ErrorDetail{
				Code:    models.ErrInvalidInput,
				Message: "invalid resolution (daily, weekly, monthly allowed)",
			},
		})
		return
	}

	// Check cache first
	cached, err := h.cache.GetHistoricalPrice(c.Request.Context(), symbol, resolution)
	if err == nil && cached != nil {
		c.JSON(http.StatusOK, models.SuccessResponse{
			Success: true,
			Data:    cached,
		})
		return
	}

	// If cache not found, fetch from provider
	historicalData, err := h.provider.GetHistoricalPrices(c.Request.Context(), symbol, resolution)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Success: false,
			Error: models.ErrorDetail{
				Code:    models.ErrServiceUnavailable,
				Message: "failed to fetch historical data",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success:   true,
		Data:      historicalData,
		Timestamp: time.Now(),
	})
}

// validateDateParameters validates the combination of date parameters
func (h *PriceHandler) validateDateParameters(date, from, to string) error {
	hasDate := date != ""
	hasFrom := from != ""
	hasTo := to != ""

	// Check for parameter conflicts
	if hasDate && (hasFrom || hasTo) {
		return fmt.Errorf("cannot use 'date' parameter with 'from'/'to' parameters")
	}

	// Check for incomplete date range
	if hasFrom && !hasTo {
		return fmt.Errorf("both 'from' and 'to' parameters are required for date range queries")
	}
	if hasTo && !hasFrom {
		return fmt.Errorf("both 'from' and 'to' parameters are required for date range queries")
	}

	// Validate date formats and future dates
	if hasDate {
		if err := h.validateSingleDate(date); err != nil {
			return err
		}
	}

	if hasFrom && hasTo {
		if err := h.validateSingleDate(from); err != nil {
			return fmt.Errorf("invalid 'from' date: %v", err)
		}
		if err := h.validateSingleDate(to); err != nil {
			return fmt.Errorf("invalid 'to' date: %v", err)
		}

		// Validate date range order
		fromDate, _ := time.Parse("2006-01-02", from)
		toDate, _ := time.Parse("2006-01-02", to)
		if fromDate.After(toDate) {
			return fmt.Errorf("'from' date must be earlier than or equal to 'to' date")
		}
	}

	return nil
}

// validateSingleDate validates a single date parameter
func (h *PriceHandler) validateSingleDate(dateStr string) error {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}

	if parsedDate.After(time.Now()) {
		return fmt.Errorf("date cannot be in the future")
	}

	return nil
}

// filterHistoricalDataByDateParams filters historical data based on date parameters
// If dateParam is provided, filters for that specific date
// If fromDate and toDate are provided, filters for the date range (inclusive)
func (h *PriceHandler) filterHistoricalDataByDateParams(data *models.SymbolHistoricalPrice, dateParam, fromDate, toDate string) *models.SymbolHistoricalPrice {
	var filteredPrices []models.ClosePrice

	// Handle single date filtering
	if dateParam != "" {
		for _, priceData := range data.HistoricalPrices {
			if priceData.Date == dateParam {
				filteredPrices = append(filteredPrices, priceData)
			}
		}
		return &models.SymbolHistoricalPrice{
			Symbol:           data.Symbol,
			Resolution:       models.ResolutionDaily, // Always daily for date-specific queries
			HistoricalPrices: filteredPrices,
		}
	}

	// Handle date range filtering
	if fromDate != "" && toDate != "" {
		fromTime, _ := time.Parse("2006-01-02", fromDate)
		toTime, _ := time.Parse("2006-01-02", toDate)

		for _, priceData := range data.HistoricalPrices {
			priceTime, err := time.Parse("2006-01-02", priceData.Date)
			if err != nil {
				continue // Skip invalid dates
			}

			// Include dates within the range (inclusive)
			if (priceTime.Equal(fromTime) || priceTime.After(fromTime)) &&
				(priceTime.Equal(toTime) || priceTime.Before(toTime)) {
				filteredPrices = append(filteredPrices, priceData)
			}
		}
	}

	return &models.SymbolHistoricalPrice{
		Symbol:           data.Symbol,
		Resolution:       models.ResolutionDaily, // Always daily for date-specific queries
		HistoricalPrices: filteredPrices,
	}
}

// checkCacheCoverage determines if the cached data covers the requested date range
func (h *PriceHandler) checkCacheCoverage(data *models.SymbolHistoricalPrice, fromDate, toDate string) string {
	if len(data.HistoricalPrices) == 0 {
		return "none"
	}

	fromTime, _ := time.Parse("2006-01-02", fromDate)
	toTime, _ := time.Parse("2006-01-02", toDate)

	// Get the latest (most recent) price record from cache (first record since sorted newest to oldest)
	lastPriceDate, err := time.Parse("2006-01-02", data.HistoricalPrices[0].Date)
	if err != nil {
		return "none"
	}

	// Compare requested date range with last price record
	// Adjust toDate to the last trading day if it falls on weekend or holiday
	adjustedToDate := h.getLastTradingDay(toTime)

	// 1. If both fromDate and toDate are later than last price record → no coverage
	if fromTime.After(lastPriceDate) && adjustedToDate.After(lastPriceDate) {
		return "none"
	}

	// 2. If fromDate is earlier than last price record and toDate is later → partial coverage
	if fromTime.Before(lastPriceDate) && adjustedToDate.After(lastPriceDate) {
		return "partial"
	}

	// 3. If last price record is later than or equal to adjusted toDate → full coverage
	if lastPriceDate.After(adjustedToDate) || lastPriceDate.Equal(adjustedToDate) {
		return "full"
	}

	// Default case: partial coverage
	return "partial"
}

// isUSMarketHoliday checks if a given date is a US stock market holiday
func (h *PriceHandler) isUSMarketHoliday(date time.Time) bool {
	month := date.Month()
	day := date.Day()

	// New Year's Day (January 1)
	if month == time.January && day == 1 {
		return true
	}

	// Martin Luther King Jr. Day (3rd Monday in January)
	if month == time.January && date.Weekday() == time.Monday {
		if day >= 15 && day <= 21 {
			return true
		}
	}

	// Presidents' Day (3rd Monday in February)
	if month == time.February && date.Weekday() == time.Monday {
		if day >= 15 && day <= 21 {
			return true
		}
	}

	// Good Friday (Friday before Easter - complex calculation, simplified for major years)
	// Memorial Day (last Monday in May)
	if month == time.May && date.Weekday() == time.Monday && day >= 25 {
		return true
	}

	// Juneteenth (June 19)
	if month == time.June && day == 19 {
		return true
	}

	// Independence Day (July 4)
	if month == time.July && day == 4 {
		return true
	}

	// Labor Day (1st Monday in September)
	if month == time.September && date.Weekday() == time.Monday && day <= 7 {
		return true
	}

	// Thanksgiving (4th Thursday in November)
	if month == time.November && date.Weekday() == time.Thursday {
		if day >= 22 && day <= 28 {
			return true
		}
	}

	// Christmas Day (December 25)
	if month == time.December && day == 25 {
		return true
	}

	return false
}

// getLastTradingDay returns the last trading day on or before the given date
func (h *PriceHandler) getLastTradingDay(date time.Time) time.Time {
	adjustedDate := date

	// Keep going back until we find a trading day
	for {
		// Skip weekends
		if adjustedDate.Weekday() == time.Saturday || adjustedDate.Weekday() == time.Sunday {
			adjustedDate = adjustedDate.AddDate(0, 0, -1)
			continue
		}

		// Skip US market holidays
		if h.isUSMarketHoliday(adjustedDate) {
			adjustedDate = adjustedDate.AddDate(0, 0, -1)
			continue
		}

		// Found a trading day
		break
	}

	return adjustedDate
}
