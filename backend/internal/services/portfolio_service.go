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
			CurrentPrice:         currentPrice,
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

// calculateAnnualizedReturnRate calculates annualized return rate based on holding period
func (s *PortfolioService) calculateAnnualizedReturnRate(transactions []models.Transaction, totalCost, marketValue float64) float64 {
	if totalCost <= 0 || len(transactions) == 0 {
		return 0
	}

	// Find the earliest buy transaction to calculate holding period
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

	// Calculate holding period in years
	holdingPeriodDays := time.Since(earliestBuyDate).Hours() / 24
	holdingPeriodYears := holdingPeriodDays / 365.25

	if holdingPeriodYears <= 0 {
		return 0
	}

	// Calculate annualized return: ((ending_value / beginning_value) ^ (1/years)) - 1
	totalReturn := marketValue / totalCost
	annualizedReturn := (math.Pow(totalReturn, 1/holdingPeriodYears) - 1) * 100

	return annualizedReturn
}
