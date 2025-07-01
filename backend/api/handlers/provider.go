package handlers

import (
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/repositories"
	"github.com/transaction-tracker/backend/internal/services"
	"gorm.io/gorm"
)

type Handlers struct {
	Transactions *TransactionsHandler
	Auth         *AuthHandler
}

// InitHandlers wires up all dependencies and returns a Handlers struct
func InitHandlers(db *gorm.DB, cfg *config.Config) *Handlers {
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)

	return &Handlers{
		Transactions: NewTransactionsHandler(transactionService),
		Auth:         NewAuthHandler(db, cfg),
	}
}
