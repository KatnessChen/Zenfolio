package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/transaction-tracker/backend/internal/models"
	"github.com/transaction-tracker/backend/internal/repositories"
	"github.com/transaction-tracker/backend/internal/utils"
)

// TransactionFilter represents filters for transaction queries
type TransactionFilter struct {
	UserID         *string
	Symbols        []string // Support multiple symbols
	TradeTypes     []string // Support multiple types
	Exchanges      []string // Support multiple exchanges
	Brokers        []string // Support multiple brokers
	Currencies     []string // Support multiple currencies
	StartDate      *time.Time
	EndDate        *time.Time
	MinAmount      *float64
	MaxAmount      *float64
	Limit          int
	Offset         int
	OrderBy        string
	OrderDirection string
}

// TransactionService handles transaction-related business logic
type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(transactionRepo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

// CreateTransactions creates multiple transactions in a batch (business logic)
func (s *TransactionService) CreateTransactions(userID string, transactions []models.Transaction) ([]models.Transaction, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	userUUID := utils.UUID{UUID: parsedUserID}
	// Set user ID for each transaction (business logic)
	for i := range transactions {
		transactions[i].UserID = userUUID
	}

	// Delegate to repository for database operations
	return s.transactionRepo.CreateMany(transactions)
}

// GetPortfolioSummary returns portfolio summary for a user
func (s *TransactionService) GetPortfolioSummary(userID string) (map[string]interface{}, error) {
	return s.transactionRepo.GetPortfolioSummaryByUserID(userID)
}

// GetSymbolHoldings returns current holdings for a user grouped by symbol
func (s *TransactionService) GetSymbolHoldings(userID string) ([]map[string]interface{}, error) {
	return s.transactionRepo.GetSymbolHoldingsByUserID(userID)
}

// GetTransactionsWithFilter retrieves transactions with advanced filtering (business logic method)
func (s *TransactionService) GetTransactionsWithFilter(filter TransactionFilter) ([]models.Transaction, error) {
	return s.transactionRepo.GetWithFilters(
		filter.UserID,
		filter.Symbols,
		filter.TradeTypes,
		filter.Exchanges,
		filter.Brokers,
		filter.Currencies,
		filter.StartDate,
		filter.EndDate,
		filter.MinAmount,
		filter.MaxAmount,
		filter.OrderBy,
		filter.OrderDirection,
		filter.Limit,
		filter.Offset,
	)
}

// CountTransactions returns the count of transactions based on filter
func (s *TransactionService) CountTransactions(filter TransactionFilter) (int64, error) {
	return s.transactionRepo.CountWithFilters(
		filter.UserID,
		filter.Symbols,
		filter.TradeTypes,
		filter.Exchanges,
		filter.Brokers,
		filter.Currencies,
		filter.StartDate,
		filter.EndDate,
		filter.MinAmount,
		filter.MaxAmount,
	)
}
