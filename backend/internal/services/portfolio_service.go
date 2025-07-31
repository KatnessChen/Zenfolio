package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/transaction-tracker/backend/internal/models"
	"github.com/transaction-tracker/backend/internal/provider"
	"github.com/transaction-tracker/backend/internal/repositories"
	"github.com/transaction-tracker/backend/internal/utils"
)

// PortfolioService handles portfolio-related business logic
type PortfolioService struct {
	transactionRepo *repositories.TransactionRepository
	priceManager    *provider.PriceServiceManager
}

// NewPortfolioService creates a new portfolio service
func NewPortfolioService(
	transactionRepo *repositories.TransactionRepository,
	priceManager *provider.PriceServiceManager,
) *PortfolioService {
	return &PortfolioService{
		transactionRepo: transactionRepo,
		priceManager:    priceManager,
	}
}

// GetSingleHoldingBasicInfo retrieves basic information for a specific stock holding
func (s *PortfolioService) GetSingleHoldingBasicInfo(ctx context.Context, userID uuid.UUID, symbol string) (*models.SingleHolding, error) {
	// Get all transactions for this user and symbol
	transactions, err := s.transactionRepo.GetByUserIDAndSymbol(userID, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for symbol %s: %w", symbol, err)
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transactions found for symbol %s", symbol)
	}

	// Calculate basic metrics from transactions
	totalQuantity, totalCost, unitCost, realizedGainLoss := s.calculateHoldingMetrics(transactions)

	// Check if user still holds this stock
	if totalQuantity <= 0 {
		return nil, fmt.Errorf("no current holdings for symbol %s", symbol)
	}

	// Get current price from PriceServiceManager
	currentPriceData, err := s.priceManager.GetCurrentPrice(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get current price for %s: %w", symbol, err)
	}

	currentPrice := currentPriceData.CurrentPrice
	marketValue := totalQuantity * currentPrice
	unrealizedGainLoss := marketValue - totalCost

	// Calculate return rates
	simpleReturnRate := s.calculateSimpleReturnRate(totalCost, marketValue)
	annualizedReturnRate := s.calculateAnnualizedReturnRate(transactions, totalCost, marketValue)

	return &models.SingleHolding{
		Symbol:               symbol,
		TotalQuantity:        utils.RoundTo4(totalQuantity),
		TotalCost:            utils.RoundTo4(totalCost),
		UnitCost:             utils.RoundTo4(unitCost),
		CurrentPrice:         currentPrice,
		MarketValue:          utils.RoundTo4(marketValue),
		SimpleReturnRate:     utils.RoundTo4(simpleReturnRate),
		AnnualizedReturnRate: utils.RoundTo4(annualizedReturnRate),
		RealizedGainLoss:     utils.RoundTo4(realizedGainLoss),
		UnrealizedGainLoss:   utils.RoundTo4(unrealizedGainLoss),
	}, nil
}

// GetAllHoldings retrieves basic information for all current holdings of a user
func (s *PortfolioService) GetAllHoldings(ctx context.Context, userID uuid.UUID) ([]models.SingleHolding, error) {
	// Get all transactions for this user
	transactions, err := s.transactionRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for user: %w", err)
	}

	if len(transactions) == 0 {
		return []models.SingleHolding{}, nil
	}

	// Group transactions by symbol
	transactionsBySymbol := make(map[string][]models.Transaction)
	for _, tx := range transactions {
		transactionsBySymbol[tx.Symbol] = append(transactionsBySymbol[tx.Symbol], tx)
	}

	var holdings []models.SingleHolding
	for symbol, symbolTransactions := range transactionsBySymbol {
		// Calculate basic metrics for this symbol
		totalQuantity, totalCost, unitCost, realizedGainLoss := s.calculateHoldingMetrics(symbolTransactions)

		// Skip if user no longer holds this stock
		if totalQuantity <= 0 {
			continue
		}

		// Get current price from PriceServiceManager
		currentPriceData, err := s.priceManager.GetCurrentPrice(ctx, symbol)
		if err != nil {
			// Log error but continue with other holdings
			fmt.Printf("Warning: failed to get current price for %s: %v\n", symbol, err)
			continue
		}

		currentPrice := currentPriceData.CurrentPrice
		marketValue := totalQuantity * currentPrice
		unrealizedGainLoss := marketValue - totalCost

		// Calculate return rates
		simpleReturnRate := s.calculateSimpleReturnRate(totalCost, marketValue)
		annualizedReturnRate := s.calculateAnnualizedReturnRate(symbolTransactions, totalCost, marketValue)

		holding := models.SingleHolding{
			Symbol:               symbol,
			TotalQuantity:        utils.RoundTo4(totalQuantity),
			TotalCost:            utils.RoundTo4(totalCost),
			UnitCost:             utils.RoundTo4(unitCost),
			CurrentPrice:         utils.RoundTo4(currentPrice),
			MarketValue:          utils.RoundTo4(marketValue),
			SimpleReturnRate:     utils.RoundTo4(simpleReturnRate),
			AnnualizedReturnRate: utils.RoundTo4(annualizedReturnRate),
			RealizedGainLoss:     utils.RoundTo4(realizedGainLoss),
			UnrealizedGainLoss:   utils.RoundTo4(unrealizedGainLoss),
		}

		holdings = append(holdings, holding)
	}

	return holdings, nil
}

// GetPortfolioSummary retrieves comprehensive portfolio summary for a user
func (s *PortfolioService) GetPortfolioSummary(ctx context.Context, userID uuid.UUID) (*models.PortfolioSummary, error) {
	// Get all current holdings
	holdings, err := s.GetAllHoldings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get holdings for portfolio summary: %w", err)
	}

	now := time.Now().UTC()

	// Calculate summary metrics
	var totalMarketValue, totalCost, totalRealizedGainLoss, totalUnrealizedGainLoss float64
	holdingsCount := len(holdings)

	// Check if user has any transactions (not just current holdings)
	allTransactions, err := s.transactionRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions for portfolio summary: %w", err)
	}
	hasTransactions := len(allTransactions) > 0

	for _, holding := range holdings {
		totalMarketValue += holding.MarketValue
		totalCost += holding.TotalCost
		totalRealizedGainLoss += holding.RealizedGainLoss
		totalUnrealizedGainLoss += holding.UnrealizedGainLoss
	}

	// Calculate total return and percentage
	totalReturn := totalUnrealizedGainLoss + totalRealizedGainLoss
	var totalReturnPercentage float64
	if totalCost > 0 {
		totalReturnPercentage = (totalReturn / totalCost) * 100
	}

	// Calculate annualized return rate (XIRR) for the whole portfolio
	var annualizedReturnRate float64
	if hasTransactions && totalCost > 0 && totalMarketValue > 0 {
		// Gather all transactions for XIRR calculation
		allTxs := append([]models.Transaction{}, allTransactions...)
		annualizedReturnRate = s.calculateAnnualizedReturnRate(allTxs, totalCost, totalMarketValue)
	}

	return &models.PortfolioSummary{
		Timestamp:             now,
		Currency:              "USD", // Default currency as per requirements
		MarketValue:           utils.RoundTo4(totalMarketValue),
		TotalCost:             utils.RoundTo4(totalCost),
		TotalReturn:           utils.RoundTo4(totalReturn),
		TotalReturnPercentage: utils.RoundTo4(totalReturnPercentage),
		HoldingsCount:         holdingsCount,
		HasTransactions:       hasTransactions,
		AnnualizedReturnRate:  utils.RoundTo4(annualizedReturnRate),
		LastUpdated:           now,
	}, nil
}

// calculateHoldingMetrics calculates total quantity, cost, unit cost, and realized gains/losses
func (s *PortfolioService) calculateHoldingMetrics(transactions []models.Transaction) (totalQuantity, totalCost, unitCost, realizedGainLoss float64) {
	var totalBoughtQuantity, totalBoughtCost float64
	var totalSoldQuantity, totalSoldRevenue float64

	for _, tx := range transactions {
		switch tx.TradeType {
		case "Buy":
			totalBoughtQuantity += tx.Quantity
			totalBoughtCost += tx.Amount
		case "Sell":
			totalSoldQuantity += tx.Quantity
			totalSoldRevenue += tx.Amount
		}
	}

	// Current holdings
	totalQuantity = totalBoughtQuantity - totalSoldQuantity

	if totalQuantity > 0 && totalBoughtQuantity > 0 {
		// Calculate weighted average cost basis for remaining shares
		avgCostPerShare := totalBoughtCost / totalBoughtQuantity
		totalCost = totalQuantity * avgCostPerShare
		unitCost = avgCostPerShare
	}

	// Realized gain/loss from sold positions
	if totalSoldQuantity > 0 && totalBoughtQuantity > 0 {
		avgCostPerShare := totalBoughtCost / totalBoughtQuantity
		soldCostBasis := totalSoldQuantity * avgCostPerShare
		realizedGainLoss = totalSoldRevenue - soldCostBasis
	}

	return totalQuantity, totalCost, unitCost, realizedGainLoss
}

// calculateSimpleReturnRate calculates simple return rate percentage
func (s *PortfolioService) calculateSimpleReturnRate(totalCost, marketValue float64) float64 {
	if totalCost <= 0 {
		return 0
	}
	return ((marketValue - totalCost) / totalCost) * 100
}

// calculateAnnualizedReturnRate calculates XIRR (Internal Rate of Return) based on cash flows
func (s *PortfolioService) calculateAnnualizedReturnRate(transactions []models.Transaction, totalCost, marketValue float64) float64 {
	if totalCost <= 0 || len(transactions) == 0 || marketValue <= 0 {
		return 0
	}

	// Build cash flows from transactions
	var cashFlows []struct {
		Amount float64
		Date   time.Time
	}

	for _, tx := range transactions {
		amount := tx.Amount
		if tx.TradeType == "Buy" {
			amount = -amount
		}
		cashFlows = append(cashFlows, struct {
			Amount float64
			Date   time.Time
		}{
			Amount: amount,
			Date:   tx.TransactionDate,
		})
	}
	cashFlows = append(cashFlows, struct {
		Amount float64
		Date   time.Time
	}{
		Amount: marketValue,
		Date:   time.Now(),
	})

	rate := utils.XIRR(cashFlows)
	if rate == 0 {
		var earliestBuyDate time.Time
		for _, tx := range transactions {
			if tx.TradeType == "Buy" {
				if earliestBuyDate.IsZero() || tx.TransactionDate.Before(earliestBuyDate) {
					earliestBuyDate = tx.TransactionDate
				}
			}
		}
		if earliestBuyDate.IsZero() {
			return 0
		}
		holdingPeriodDays := time.Since(earliestBuyDate).Hours() / 24
		holdingPeriodYears := holdingPeriodDays / 365.25
		if holdingPeriodYears <= 0 {
			return 0
		}
		totalReturn := marketValue / totalCost
		return (math.Pow(totalReturn, 1/holdingPeriodYears) - 1) * 100
	}
	return rate * 100
}

// GetHistoricalPortfolioTotalValue calculates portfolio total value over time
func (s *PortfolioService) GetHistoricalPortfolioTotalValue(ctx context.Context, userID uuid.UUID, timeframe models.TimeFrame) (*models.HistoricalTotalValueResponse, error) {
	// Calculate time range based on timeframe
	endTime := time.Now()
	startTime, err := s.calculateStartTime(endTime, timeframe)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate start time: %w", err)
	}

	// Determine granularity
	defaultGranularity := s.determineDefaultGranularity(timeframe)
	var granularity = &defaultGranularity

	// Get all transactions for the user up to end time (needed for correct portfolio calculation)
	allTransactions, err := s.transactionRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// For ALL timeframe, use first transaction date as start time
	if timeframe == models.TimeFrameALL && len(allTransactions) > 0 {
		startTime = allTransactions[0].TransactionDate
	}

	if len(allTransactions) == 0 {
		return &models.HistoricalTotalValueResponse{
			TimeFrame:   timeframe,
			Granularity: *granularity,
			Period: struct {
				StartDate time.Time `json:"start_date"`
				EndDate   time.Time `json:"end_date"`
			}{
				StartDate: startTime,
				EndDate:   endTime,
			},
			DataPoints: []models.TotalValueDataPoint{},
			Summary: models.TotalValueTrendSummary{
				Change:        0,
				ChangePercent: 0,
				Volatility:    0,
				MaxValue:      0,
				MinValue:      0,
			},
		}, nil
	}

	// Generate time points based on granularity
	timePoints := s.generateTimePoints(startTime, endTime, *granularity)

	// Calculate total value for each time point
	dataPoints := make([]models.TotalValueDataPoint, 0, len(timePoints))
	var previousValue float64

	for i, timePoint := range timePoints {
		totalValue, err := s.calculateTotalValueAtTime(ctx, allTransactions, timePoint)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate total value at %v: %w", timePoint, err)
		}

		// Calculate day change
		dayChange := 0.0
		dayChangePercent := 0.0
		if i > 0 && previousValue > 0 {
			dayChange = totalValue - previousValue
			dayChangePercent = (dayChange / previousValue) * 100
		}

		dataPoints = append(dataPoints, models.TotalValueDataPoint{
			Timestamp:        timePoint,
			TotalValue:       totalValue,
			DayChange:        dayChange,
			DayChangePercent: dayChangePercent,
		})

		previousValue = totalValue
	}

	// Calculate summary statistics
	summary := s.calculateSummaryStatistics(dataPoints)

	return &models.HistoricalTotalValueResponse{
		TimeFrame:   timeframe,
		Granularity: *granularity,
		Period: struct {
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
		}{
			StartDate: startTime,
			EndDate:   endTime,
		},
		DataPoints: dataPoints,
		Summary:    summary,
	}, nil
}

// calculateStartTime determines the start time based on timeframe
func (s *PortfolioService) calculateStartTime(endTime time.Time, timeframe models.TimeFrame) (time.Time, error) {
	switch timeframe {
	case models.TimeFrame1D:
		return endTime.AddDate(0, 0, -1), nil
	case models.TimeFrame1W:
		return endTime.AddDate(0, 0, -7), nil
	case models.TimeFrame1M:
		return endTime.AddDate(0, -1, 0), nil
	case models.TimeFrame3M:
		return endTime.AddDate(0, -3, 0), nil
	case models.TimeFrame6M:
		return endTime.AddDate(0, -6, 0), nil
	case models.TimeFrameYTD:
		return time.Date(endTime.Year(), 1, 1, 0, 0, 0, 0, endTime.Location()), nil
	case models.TimeFrame1Y:
		return endTime.AddDate(-1, 0, 0), nil
	case models.TimeFrame5Y:
		return endTime.AddDate(-5, 0, 0), nil
	case models.TimeFrameALL:
		// For ALL, we'll start from the first transaction date
		// This will be handled separately in the main function
		return time.Time{}, nil
	default:
		return time.Time{}, fmt.Errorf("unsupported timeframe: %s", timeframe)
	}
}

// determineDefaultGranularity determines appropriate granularity for timeframe
func (s *PortfolioService) determineDefaultGranularity(timeframe models.TimeFrame) models.Granularity {
	switch timeframe {
	case models.TimeFrame1W:
		return models.GranularityDaily
	case models.TimeFrame1M:
		return models.GranularityWeekly
	case models.TimeFrame3M, models.TimeFrame6M:
		return models.GranularityMonthly
	case models.TimeFrameYTD, models.TimeFrame1Y, models.TimeFrame5Y, models.TimeFrameALL:
		return models.GranularityMonthly
	default:
		return models.GranularityDaily
	}
}

// generateTimePoints creates time points based on granularity
func (s *PortfolioService) generateTimePoints(startTime, endTime time.Time, granularity models.Granularity) []time.Time {
	var timePoints []time.Time
	current := startTime

	for current.Before(endTime) || current.Equal(endTime) {
		timePoints = append(timePoints, current)

		switch granularity {
		case models.GranularityHourly:
			current = current.Add(time.Hour)
		case models.GranularityDaily:
			current = current.AddDate(0, 0, 1)
		case models.GranularityWeekly:
			current = current.AddDate(0, 0, 7)
		case models.GranularityMonthly:
			current = current.AddDate(0, 1, 0)
		}
	}

	// Ensure endTime is included for Weekly and Monthly if not already present
	if (granularity == models.GranularityWeekly || granularity == models.GranularityMonthly) && len(timePoints) > 0 {
		last := timePoints[len(timePoints)-1]
		if last.Before(endTime) {
			timePoints = append(timePoints, endTime)
		}
	}

	return timePoints
}

// calculateTotalValueAtTime calculates portfolio total value at a specific time
func (s *PortfolioService) calculateTotalValueAtTime(ctx context.Context, transactions []models.Transaction, targetTime time.Time) (float64, error) {
	// Group transactions by symbol and calculate holdings at target time
	holdings := make(map[string]float64)

	for _, transaction := range transactions {
		if transaction.TransactionDate.After(targetTime) {
			continue // Skip future transactions
		}

		symbol := transaction.Symbol
		quantity := transaction.Quantity

		if transaction.TradeType == "Sell" {
			quantity = -quantity
		}

		holdings[symbol] += quantity
	}

	// Calculate total market value using historical prices at target time
	totalValue := 0.0
	targetDateStr := targetTime.Format("2006-01-02")

	for symbol, quantity := range holdings {
		if quantity <= 0 {
			continue // Skip if no holdings
		}

		// Get historical price for the symbol at target date
		historicalPrice, err := s.priceManager.GetHistoricalPriceAtDate(ctx, symbol, targetDateStr)
		if err != nil {
			// If we can't get historical price, fallback to current price as last resort
			priceData, fallbackErr := s.priceManager.GetCurrentPrice(ctx, symbol)
			if fallbackErr != nil {
				// If both historical and current price fail, skip this holding
				continue
			}
			totalValue += quantity * priceData.CurrentPrice
			continue
		}

		// Find the price for the exact date
		var priceAtDate float64
		found := false

		for _, pricePoint := range historicalPrice.HistoricalPrices {
			if pricePoint.Date == targetDateStr {
				priceAtDate = pricePoint.Price
				found = true
				break
			}
		}

		if !found {
			// If exact date not found, fallback to current price
			priceData, fallbackErr := s.priceManager.GetCurrentPrice(ctx, symbol)
			if fallbackErr != nil {
				continue
			}
			priceAtDate = priceData.CurrentPrice
		}

		totalValue += quantity * priceAtDate
	}

	return totalValue, nil
}

// calculateSummaryStatistics calculates summary statistics for the data points
func (s *PortfolioService) calculateSummaryStatistics(dataPoints []models.TotalValueDataPoint) models.TotalValueTrendSummary {
	if len(dataPoints) == 0 {
		return models.TotalValueTrendSummary{}
	}

	firstValue := dataPoints[0].TotalValue
	lastValue := dataPoints[len(dataPoints)-1].TotalValue

	change := lastValue - firstValue
	changePercent := 0.0
	if firstValue > 0 {
		changePercent = (change / firstValue) * 100
	}

	minValue := dataPoints[0].TotalValue
	maxValue := dataPoints[0].TotalValue

	var returns []float64
	for i := 1; i < len(dataPoints); i++ {
		if dataPoints[i-1].TotalValue > 0 {
			dailyReturn := (dataPoints[i].TotalValue - dataPoints[i-1].TotalValue) / dataPoints[i-1].TotalValue
			returns = append(returns, dailyReturn)
		}
		if dataPoints[i].TotalValue > maxValue {
			maxValue = dataPoints[i].TotalValue
		}
		if dataPoints[i].TotalValue < minValue {
			minValue = dataPoints[i].TotalValue
		}
	}

	volatility := utils.StandardDeviation(returns) * 100

	return models.TotalValueTrendSummary{
		Change:        change,
		ChangePercent: changePercent,
		Volatility:    volatility,
		MaxValue:      maxValue,
		MinValue:      minValue,
	}
}
