package handlers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/transaction-tracker/price_service/internal/models"
)

func TestValidateDateParameters(t *testing.T) {
	handler := &PriceHandler{}

	tests := []struct {
		name        string
		date        string
		from        string
		to          string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid single date",
			date:        "2025-01-01",
			from:        "",
			to:          "",
			expectError: false,
		},
		{
			name:        "Valid date range",
			date:        "",
			from:        "2025-01-01",
			to:          "2025-01-10",
			expectError: false,
		},
		{
			name:        "No parameters (valid for resolution query)",
			date:        "",
			from:        "",
			to:          "",
			expectError: false,
		},
		{
			name:        "Conflicting parameters",
			date:        "2025-01-01",
			from:        "2025-01-01",
			to:          "2025-01-10",
			expectError: true,
			errorMsg:    "cannot use 'date' parameter with 'from'/'to' parameters",
		},
		{
			name:        "Incomplete date range - only from",
			date:        "",
			from:        "2025-01-01",
			to:          "",
			expectError: true,
			errorMsg:    "both 'from' and 'to' parameters are required for date range queries",
		},
		{
			name:        "Incomplete date range - only to",
			date:        "",
			from:        "",
			to:          "2025-01-10",
			expectError: true,
			errorMsg:    "both 'from' and 'to' parameters are required for date range queries",
		},
		{
			name:        "Invalid date range order",
			date:        "",
			from:        "2025-01-10",
			to:          "2025-01-01",
			expectError: true,
			errorMsg:    "'from' date must be earlier than or equal to 'to' date",
		},
		{
			name:        "Future date",
			date:        "2030-01-01",
			from:        "",
			to:          "",
			expectError: true,
			errorMsg:    "date cannot be in the future",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateDateParameters(tt.date, tt.from, tt.to)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSingleDate(t *testing.T) {
	handler := &PriceHandler{}

	tests := []struct {
		name        string
		dateStr     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid date",
			dateStr:     "2025-01-01",
			expectError: false,
		},
		{
			name:        "Invalid format - wrong separator",
			dateStr:     "2025/01/01",
			expectError: true,
			errorMsg:    "invalid date format, use YYYY-MM-DD",
		},
		{
			name:        "Invalid format - wrong order",
			dateStr:     "01-01-2025",
			expectError: true,
			errorMsg:    "invalid date format, use YYYY-MM-DD",
		},
		{
			name:        "Invalid date - June 31st",
			dateStr:     "2025-06-31",
			expectError: true,
			errorMsg:    "invalid date format, use YYYY-MM-DD",
		},
		{
			name:        "Invalid date - February 30th",
			dateStr:     "2025-02-30",
			expectError: true,
			errorMsg:    "invalid date format, use YYYY-MM-DD",
		},
		{
			name:        "Future date",
			dateStr:     "2030-01-01",
			expectError: true,
			errorMsg:    "date cannot be in the future",
		},
		{
			name:        "Leap year valid",
			dateStr:     "2024-02-29",
			expectError: false,
		},
		{
			name:        "Non-leap year invalid",
			dateStr:     "2025-02-29",
			expectError: true,
			errorMsg:    "invalid date format, use YYYY-MM-DD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateSingleDate(tt.dateStr)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFilterHistoricalDataByDateParams(t *testing.T) {
	handler := &PriceHandler{}

	// Sample historical data
	data := &models.SymbolHistoricalPrice{
		Symbol:     "NVDA",
		Resolution: models.ResolutionDaily,
		HistoricalPrices: []models.ClosePrice{
			{Date: "2025-01-15", Price: 150.00},
			{Date: "2025-01-10", Price: 140.00},
			{Date: "2025-01-05", Price: 130.00},
			{Date: "2024-12-31", Price: 120.00},
			{Date: "2024-12-25", Price: 110.00},
		},
	}

	tests := []struct {
		name          string
		dateParam     string
		fromDate      string
		toDate        string
		expectedCount int
		expectedDates []string
	}{
		{
			name:          "Single date filter",
			dateParam:     "2025-01-10",
			fromDate:      "",
			toDate:        "",
			expectedCount: 1,
			expectedDates: []string{"2025-01-10"},
		},
		{
			name:          "Single date not found",
			dateParam:     "2025-01-12",
			fromDate:      "",
			toDate:        "",
			expectedCount: 0,
			expectedDates: []string{},
		},
		{
			name:          "Date range covers multiple dates",
			dateParam:     "",
			fromDate:      "2025-01-01",
			toDate:        "2025-01-12",
			expectedCount: 2,
			expectedDates: []string{"2025-01-10", "2025-01-05"},
		},
		{
			name:          "Date range covers single date",
			dateParam:     "",
			fromDate:      "2025-01-05",
			toDate:        "2025-01-05",
			expectedCount: 1,
			expectedDates: []string{"2025-01-05"},
		},
		{
			name:          "Date range covers no dates",
			dateParam:     "",
			fromDate:      "2025-02-01",
			toDate:        "2025-02-10",
			expectedCount: 0,
			expectedDates: []string{},
		},
		{
			name:          "Date range extends beyond available data",
			dateParam:     "",
			fromDate:      "2024-12-20",
			toDate:        "2025-01-20",
			expectedCount: 5,
			expectedDates: []string{"2025-01-15", "2025-01-10", "2025-01-05", "2024-12-31", "2024-12-25"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.filterHistoricalDataByDateParams(data, tt.dateParam, tt.fromDate, tt.toDate)

			assert.Equal(t, tt.expectedCount, len(result.HistoricalPrices))
			assert.Equal(t, "NVDA", result.Symbol)
			assert.Equal(t, models.ResolutionDaily, result.Resolution)

			for i, expectedDate := range tt.expectedDates {
				assert.Equal(t, expectedDate, result.HistoricalPrices[i].Date)
			}
		})
	}
}

func TestIsUSMarketHoliday(t *testing.T) {
	handler := &PriceHandler{}

	tests := []struct {
		name     string
		date     string
		expected bool
	}{
		{
			name:     "New Year's Day",
			date:     "2025-01-01",
			expected: true,
		},
		{
			name:     "MLK Day - 3rd Monday in January",
			date:     "2025-01-20",
			expected: true,
		},
		{
			name:     "Presidents' Day - 3rd Monday in February",
			date:     "2025-02-17",
			expected: true,
		},
		{
			name:     "Memorial Day - last Monday in May",
			date:     "2025-05-26",
			expected: true,
		},
		{
			name:     "Juneteenth",
			date:     "2025-06-19",
			expected: true,
		},
		{
			name:     "Independence Day",
			date:     "2025-07-04",
			expected: true,
		},
		{
			name:     "Labor Day - 1st Monday in September",
			date:     "2025-09-01",
			expected: true,
		},
		{
			name:     "Thanksgiving - 4th Thursday in November",
			date:     "2025-11-27",
			expected: true,
		},
		{
			name:     "Christmas Day",
			date:     "2025-12-25",
			expected: true,
		},
		{
			name:     "Regular trading day",
			date:     "2025-07-23",
			expected: false,
		},
		{
			name:     "Weekend but not holiday",
			date:     "2025-07-26", // Saturday
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, _ := time.Parse("2006-01-02", tt.date)
			result := handler.isUSMarketHoliday(date)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetLastTradingDay(t *testing.T) {
	handler := &PriceHandler{}

	tests := []struct {
		name     string
		date     string
		expected string
	}{
		{
			name:     "Regular weekday",
			date:     "2025-07-23", // Wednesday
			expected: "2025-07-23",
		},
		{
			name:     "Saturday - should go to Friday",
			date:     "2025-07-26", // Saturday
			expected: "2025-07-25", // Friday
		},
		{
			name:     "Sunday - should go to Friday",
			date:     "2025-07-27", // Sunday
			expected: "2025-07-25", // Friday
		},
		{
			name:     "Holiday (July 4th Friday) - should go to Thursday",
			date:     "2025-07-04", // Friday, Independence Day
			expected: "2025-07-03", // Thursday
		},
		{
			name:     "Weekend after holiday",
			date:     "2025-07-06", // Sunday after July 4th
			expected: "2025-07-03", // Thursday (skipping Friday holiday)
		},
		{
			name:     "Christmas Day - should go to previous trading day",
			date:     "2025-12-25", // Thursday, Christmas
			expected: "2025-12-24", // Wednesday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, _ := time.Parse("2006-01-02", tt.date)
			result := handler.getLastTradingDay(date)
			assert.Equal(t, tt.expected, result.Format("2006-01-02"))
		})
	}
}

func TestCheckCacheCoverage(t *testing.T) {
	handler := &PriceHandler{}

	tests := []struct {
		name     string
		data     *models.SymbolHistoricalPrice
		fromDate string
		toDate   string
		expected string
	}{
		{
			name: "Empty cache",
			data: &models.SymbolHistoricalPrice{
				HistoricalPrices: []models.ClosePrice{},
			},
			fromDate: "2025-07-01",
			toDate:   "2025-07-10",
			expected: "none",
		},
		{
			name: "Full coverage - cache has recent data",
			data: &models.SymbolHistoricalPrice{
				HistoricalPrices: []models.ClosePrice{
					{Date: "2025-07-23", Price: 150.00}, // Most recent
					{Date: "2025-07-22", Price: 149.00},
				},
			},
			fromDate: "2025-07-20",
			toDate:   "2025-07-23",
			expected: "full",
		},
		{
			name: "No coverage - request is after cache",
			data: &models.SymbolHistoricalPrice{
				HistoricalPrices: []models.ClosePrice{
					{Date: "2025-07-20", Price: 150.00}, // Most recent
				},
			},
			fromDate: "2025-07-25",
			toDate:   "2025-07-30",
			expected: "none",
		},
		{
			name: "Partial coverage - request spans across cache boundary",
			data: &models.SymbolHistoricalPrice{
				HistoricalPrices: []models.ClosePrice{
					{Date: "2025-07-23", Price: 150.00}, // Most recent
				},
			},
			fromDate: "2025-07-20",
			toDate:   "2025-07-25",
			expected: "partial",
		},
		{
			name: "Full coverage with weekend adjustment",
			data: &models.SymbolHistoricalPrice{
				HistoricalPrices: []models.ClosePrice{
					{Date: "2025-07-25", Price: 150.00}, // Friday
				},
			},
			fromDate: "2025-07-20",
			toDate:   "2025-07-27", // Sunday - should adjust to Friday
			expected: "full",
		},
		{
			name: "Full coverage with holiday adjustment",
			data: &models.SymbolHistoricalPrice{
				HistoricalPrices: []models.ClosePrice{
					{Date: "2025-07-03", Price: 150.00}, // Thursday before July 4th
				},
			},
			fromDate: "2025-07-01",
			toDate:   "2025-07-06", // Sunday after July 4th - should adjust to July 3rd
			expected: "full",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.checkCacheCoverage(tt.data, tt.fromDate, tt.toDate)
			assert.Equal(t, tt.expected, result)
		})
	}
}
