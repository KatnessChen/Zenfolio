package repositories

import (
	"fmt"
	"time"

	"github.com/transaction-tracker/backend/internal/models"
	"github.com/transaction-tracker/backend/internal/utils"
	"gorm.io/gorm"
)

// TransactionRepository handles transaction database operations
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a single transaction
func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

// CreateMany creates multiple transactions in a single database transaction
func (r *TransactionRepository) CreateMany(transactions []models.Transaction) ([]models.Transaction, error) {
	// Start a database transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Ensure rollback on any error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var createdTransactions []models.Transaction

	for i, transaction := range transactions {
		// Create transaction
		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create transaction %d: %w", i+1, err)
		}
		createdTransactions = append(createdTransactions, transaction)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return createdTransactions, nil
}

// GetByID retrieves a transaction by transaction_id (UUID)
func (r *TransactionRepository) GetByID(id string) (*models.Transaction, error) {
	var transaction models.Transaction
	uuidBytes, err := utils.ParseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction_id: %w", err)
	}
	err = r.db.Preload("User").Where("transaction_id = ?", uuidBytes.UUID[:]).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetByUserID retrieves all transactions for a user by user_id (UUID)
func (r *TransactionRepository) GetByUserID(userID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	uuidBytes, err := utils.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}
	err = r.db.Where("user_id = ?", uuidBytes.UUID[:]).Preload("User").Find(&transactions).Error
	return transactions, err
}

// UpdateByID updates a transaction by transaction_id (UUID)
func (r *TransactionRepository) UpdateByID(id string, updates map[string]interface{}) error {
	uuidBytes, err := utils.ParseUUID(id)
	if err != nil {
		return fmt.Errorf("invalid transaction_id: %w", err)
	}
	return r.db.Model(&models.Transaction{}).Where("transaction_id = ?", uuidBytes.UUID[:]).Updates(updates).Error
}

// DeleteByID soft deletes a transaction by transaction_id (UUID)
func (r *TransactionRepository) DeleteByID(id string) error {
	uuidBytes, err := utils.ParseUUID(id)
	if err != nil {
		return fmt.Errorf("invalid transaction_id: %w", err)
	}
	return r.db.Where("transaction_id = ?", uuidBytes.UUID[:]).Delete(&models.Transaction{}).Error
}

// GetWithFilters retrieves transactions with advanced filtering
func (r *TransactionRepository) GetWithFilters(userID *string, symbols []string, types []string, exchanges []string, brokers []string, currencies []string,
	startDate *time.Time, endDate *time.Time, minAmount *float64, maxAmount *float64,
	orderBy string, orderDirection string, limit int, offset int) ([]models.Transaction, error) {

	var transactions []models.Transaction
	query := r.db.Model(&models.Transaction{})

	// Apply filters
	if userID != nil {
		uuidBytes, err := utils.ParseUUID(*userID)
		if err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		query = query.Where("user_id = ?", uuidBytes.UUID[:])
	}
	if len(symbols) > 0 {
		query = query.Where("symbol IN ?", symbols)
	}
	if len(types) > 0 {
		query = query.Where("trade_type IN ?", types)
	}
	if len(exchanges) > 0 {
		query = query.Where("exchange IN ?", exchanges)
	}
	if len(brokers) > 0 {
		query = query.Where("broker IN ?", brokers)
	}
	if len(currencies) > 0 {
		query = query.Where("currency IN ?", currencies)
	}
	if startDate != nil {
		query = query.Where("transaction_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("transaction_date <= ?", *endDate)
	}
	if minAmount != nil {
		query = query.Where("amount >= ?", *minAmount)
	}
	if maxAmount != nil {
		query = query.Where("amount <= ?", *maxAmount)
	}

	// Apply ordering
	if orderBy == "" {
		orderBy = "transaction_date"
	}
	if orderDirection == "" {
		orderDirection = "DESC"
	}
	query = query.Order(fmt.Sprintf("%s %s", orderBy, orderDirection))

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	// Preload user data
	query = query.Preload("User")

	if err := query.Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get filtered transactions: %w", err)
	}

	return transactions, nil
}

// CountWithFilters returns the count of transactions based on filters
func (r *TransactionRepository) CountWithFilters(userID *string, symbols []string, types []string, exchanges []string, brokers []string, currencies []string,
	startDate *time.Time, endDate *time.Time, minAmount *float64, maxAmount *float64) (int64, error) {

	var count int64
	query := r.db.Model(&models.Transaction{})

	// Apply filters
	if userID != nil {
		uuidBytes, err := utils.ParseUUID(*userID)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id: %w", err)
		}
		query = query.Where("user_id = ?", uuidBytes.UUID[:])
	}
	if len(symbols) > 0 {
		query = query.Where("symbol IN ?", symbols)
	}
	if len(types) > 0 {
		query = query.Where("trade_type IN ?", types)
	}
	if len(exchanges) > 0 {
		query = query.Where("exchange IN ?", exchanges)
	}
	if len(brokers) > 0 {
		query = query.Where("broker IN ?", brokers)
	}
	if len(currencies) > 0 {
		query = query.Where("currency IN ?", currencies)
	}
	if startDate != nil {
		query = query.Where("transaction_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("transaction_date <= ?", *endDate)
	}
	if minAmount != nil {
		query = query.Where("amount >= ?", *minAmount)
	}
	if maxAmount != nil {
		query = query.Where("amount <= ?", *maxAmount)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count filtered transactions: %w", err)
	}

	return count, nil
}

// GetPortfolioSummaryByUserID returns portfolio summary for a user
func (r *TransactionRepository) GetPortfolioSummaryByUserID(userID string) (map[string]interface{}, error) {
	var result struct {
		TotalTransactions int64   `json:"total_transactions"`
		TotalBuyAmount    float64 `json:"total_buy_amount"`
		TotalSellAmount   float64 `json:"total_sell_amount"`
		UniqueSymbols     int64   `json:"unique_symbols"`
	}

	uuidBytes, err := utils.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// Count total transactions
	if err := r.db.Model(&models.Transaction{}).Where("user_id = ?", uuidBytes.UUID[:]).Count(&result.TotalTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Sum buy amounts
	if err := r.db.Model(&models.Transaction{}).Where("user_id = ? AND trade_type = ?", uuidBytes.UUID[:], "Buy").Select("COALESCE(SUM(amount), 0)").Scan(&result.TotalBuyAmount).Error; err != nil {
		return nil, fmt.Errorf("failed to sum buy amounts: %w", err)
	}

	// Sum sell amounts
	if err := r.db.Model(&models.Transaction{}).Where("user_id = ? AND trade_type = ?", uuidBytes.UUID[:], "Sell").Select("COALESCE(SUM(amount), 0)").Scan(&result.TotalSellAmount).Error; err != nil {
		return nil, fmt.Errorf("failed to sum sell amounts: %w", err)
	}

	// Count unique symbols
	if err := r.db.Model(&models.Transaction{}).Where("user_id = ?", uuidBytes.UUID[:]).Distinct("symbol").Count(&result.UniqueSymbols).Error; err != nil {
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

// GetSymbolHoldingsByUserID returns current holdings for a user grouped by symbol
func (r *TransactionRepository) GetSymbolHoldingsByUserID(userID string) ([]map[string]interface{}, error) {
	var holdings []struct {
		Symbol       string  `json:"symbol"`
		TotalBought  float64 `json:"total_bought"`
		TotalSold    float64 `json:"total_sold"`
		NetQuantity  float64 `json:"net_quantity"`
		AvgBuyPrice  float64 `json:"avg_buy_price"`
		AvgSellPrice float64 `json:"avg_sell_price"`
	}

	uuidBytes, err := utils.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	query := `
		SELECT
			symbol,
			COALESCE(SUM(CASE WHEN trade_type = 'Buy' THEN quantity ELSE 0 END), 0) as total_bought,
			COALESCE(SUM(CASE WHEN trade_type = 'Sell' THEN quantity ELSE 0 END), 0) as total_sold,
			COALESCE(SUM(CASE WHEN trade_type = 'Buy' THEN quantity WHEN trade_type = 'Sell' THEN -quantity ELSE 0 END), 0) as net_quantity,
			COALESCE(AVG(CASE WHEN trade_type = 'Buy' THEN price ELSE NULL END), 0) as avg_buy_price,
			COALESCE(AVG(CASE WHEN trade_type = 'Sell' THEN price ELSE NULL END), 0) as avg_sell_price
		FROM transactions
		WHERE user_id = ? AND deleted_at IS NULL
		GROUP BY symbol
		HAVING net_quantity != 0
		ORDER BY symbol
	`

	if err := r.db.Raw(query, uuidBytes.UUID[:]).Scan(&holdings).Error; err != nil {
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
