package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/transaction-tracker/backend/internal/models"
	"gorm.io/gorm"
)

// TransactionFilter represents filters for transaction queries
type TransactionFilter struct {
	UserID         *uint
	Symbol         *string
	Type           *string
	Status         *string
	Broker         *string
	StartDate      *time.Time
	EndDate        *time.Time
	MinAmount      *float64
	MaxAmount      *float64
	ExtractedFrom  *string
	Limit          int
	Offset         int
	OrderBy        string
	OrderDirection string
}

// TransactionService handles transaction-related database operations
type TransactionService struct {
	db *gorm.DB
}

// NewTransactionService creates a new transaction service
func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

// CreateTransaction creates a new transaction
func (s *TransactionService) CreateTransaction(transaction *models.Transaction) error {
	if err := s.db.Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

// CreateTransactions creates multiple transactions in a batch
func (s *TransactionService) CreateTransactions(transactions []models.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	if err := s.db.CreateInBatches(transactions, 100).Error; err != nil {
		return fmt.Errorf("failed to create transactions: %w", err)
	}
	return nil
}

// GetTransactionByID retrieves a transaction by ID
func (s *TransactionService) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := s.db.Preload("User").First(&transaction, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return &transaction, nil
}

// GetTransactionsByUser retrieves transactions for a specific user
func (s *TransactionService) GetTransactionsByUser(userID uint, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := s.db.Where("user_id = ?", userID).Order("transaction_date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}

	return transactions, nil
}

// GetTransactionsWithFilter retrieves transactions with advanced filtering
func (s *TransactionService) GetTransactionsWithFilter(filter TransactionFilter) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := s.db.Model(&models.Transaction{})

	// Apply filters
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Symbol != nil {
		query = query.Where("symbol = ?", *filter.Symbol)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Broker != nil {
		query = query.Where("broker = ?", *filter.Broker)
	}
	if filter.StartDate != nil {
		query = query.Where("transaction_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("transaction_date <= ?", *filter.EndDate)
	}
	if filter.MinAmount != nil {
		query = query.Where("amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		query = query.Where("amount <= ?", *filter.MaxAmount)
	}
	if filter.ExtractedFrom != nil {
		query = query.Where("extracted_from = ?", *filter.ExtractedFrom)
	}

	// Apply ordering
	orderBy := "transaction_date"
	orderDirection := "DESC"
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
	}
	if filter.OrderDirection != "" {
		orderDirection = filter.OrderDirection
	}
	query = query.Order(fmt.Sprintf("%s %s", orderBy, orderDirection))

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Preload user data
	query = query.Preload("User")

	if err := query.Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get filtered transactions: %w", err)
	}

	return transactions, nil
}

// UpdateTransaction updates a transaction
func (s *TransactionService) UpdateTransaction(transaction *models.Transaction) error {
	if err := s.db.Save(transaction).Error; err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

// DeleteTransaction soft deletes a transaction
func (s *TransactionService) DeleteTransaction(id uint) error {
	if err := s.db.Delete(&models.Transaction{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

// GetTransactionsBySymbol retrieves transactions for a specific symbol
func (s *TransactionService) GetTransactionsBySymbol(symbol string, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := s.db.Where("symbol = ?", symbol).Order("transaction_date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Preload("User").Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get symbol transactions: %w", err)
	}

	return transactions, nil
}

// GetTransactionsByDateRange retrieves transactions within a date range
func (s *TransactionService) GetTransactionsByDateRange(startDate, endDate time.Time, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := s.db.Where("transaction_date BETWEEN ? AND ?", startDate, endDate).Order("transaction_date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Preload("User").Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions by date range: %w", err)
	}

	return transactions, nil
}

// GetPortfolioSummary returns portfolio summary for a user
func (s *TransactionService) GetPortfolioSummary(userID uint) (map[string]interface{}, error) {
	var result struct {
		TotalTransactions int64   `json:"total_transactions"`
		TotalBuyAmount    float64 `json:"total_buy_amount"`
		TotalSellAmount   float64 `json:"total_sell_amount"`
		UniqueSymbols     int64   `json:"unique_symbols"`
	}

	// Count total transactions
	if err := s.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&result.TotalTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Sum buy amounts
	if err := s.db.Model(&models.Transaction{}).Where("user_id = ? AND type = ?", userID, "buy").Select("COALESCE(SUM(amount), 0)").Scan(&result.TotalBuyAmount).Error; err != nil {
		return nil, fmt.Errorf("failed to sum buy amounts: %w", err)
	}

	// Sum sell amounts
	if err := s.db.Model(&models.Transaction{}).Where("user_id = ? AND type = ?", userID, "sell").Select("COALESCE(SUM(amount), 0)").Scan(&result.TotalSellAmount).Error; err != nil {
		return nil, fmt.Errorf("failed to sum sell amounts: %w", err)
	}

	// Count unique symbols
	if err := s.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Distinct("symbol").Count(&result.UniqueSymbols).Error; err != nil {
		return nil, fmt.Errorf("failed to count unique symbols: %w", err)
	}

	return map[string]interface{}{
		"total_transactions": result.TotalTransactions,
		"total_buy_amount":   result.TotalBuyAmount,
		"total_sell_amount":  result.TotalSellAmount,
		"unique_symbols":     result.UniqueSymbols,
		"net_amount":         result.TotalSellAmount - result.TotalBuyAmount,
	}, nil
}

// GetSymbolHoldings returns current holdings for a user grouped by symbol
func (s *TransactionService) GetSymbolHoldings(userID uint) ([]map[string]interface{}, error) {
	var holdings []struct {
		Symbol       string  `json:"symbol"`
		TotalBought  float64 `json:"total_bought"`
		TotalSold    float64 `json:"total_sold"`
		NetQuantity  float64 `json:"net_quantity"`
		AvgBuyPrice  float64 `json:"avg_buy_price"`
		AvgSellPrice float64 `json:"avg_sell_price"`
	}

	query := `
		SELECT 
			symbol,
			COALESCE(SUM(CASE WHEN type = 'buy' THEN quantity ELSE 0 END), 0) as total_bought,
			COALESCE(SUM(CASE WHEN type = 'sell' THEN quantity ELSE 0 END), 0) as total_sold,
			COALESCE(SUM(CASE WHEN type = 'buy' THEN quantity WHEN type = 'sell' THEN -quantity ELSE 0 END), 0) as net_quantity,
			COALESCE(AVG(CASE WHEN type = 'buy' THEN price ELSE NULL END), 0) as avg_buy_price,
			COALESCE(AVG(CASE WHEN type = 'sell' THEN price ELSE NULL END), 0) as avg_sell_price
		FROM transactions 
		WHERE user_id = ? AND deleted_at IS NULL
		GROUP BY symbol
		HAVING net_quantity != 0
		ORDER BY symbol
	`

	if err := s.db.Raw(query, userID).Scan(&holdings).Error; err != nil {
		return nil, fmt.Errorf("failed to get symbol holdings: %w", err)
	}

	result := make([]map[string]interface{}, len(holdings))
	for i, holding := range holdings {
		result[i] = map[string]interface{}{
			"symbol":         holding.Symbol,
			"total_bought":   holding.TotalBought,
			"total_sold":     holding.TotalSold,
			"net_quantity":   holding.NetQuantity,
			"avg_buy_price":  holding.AvgBuyPrice,
			"avg_sell_price": holding.AvgSellPrice,
		}
	}

	return result, nil
}

// CountTransactions returns the count of transactions based on filter
func (s *TransactionService) CountTransactions(filter TransactionFilter) (int64, error) {
	var count int64
	query := s.db.Model(&models.Transaction{})

	// Apply same filters as GetTransactionsWithFilter
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Symbol != nil {
		query = query.Where("symbol = ?", *filter.Symbol)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Broker != nil {
		query = query.Where("broker = ?", *filter.Broker)
	}
	if filter.StartDate != nil {
		query = query.Where("transaction_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("transaction_date <= ?", *filter.EndDate)
	}
	if filter.MinAmount != nil {
		query = query.Where("amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		query = query.Where("amount <= ?", *filter.MaxAmount)
	}
	if filter.ExtractedFrom != nil {
		query = query.Where("extracted_from = ?", *filter.ExtractedFrom)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	return count, nil
}
