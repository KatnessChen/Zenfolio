package handlers

import (
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/ai"
	"github.com/transaction-tracker/backend/internal/provider"
	"github.com/transaction-tracker/backend/internal/repositories"
	"github.com/transaction-tracker/backend/internal/services"
	"gorm.io/gorm"
)

type Handlers struct {
	Transactions               *TransactionsHandler
	ExtractTransactionsHandler *ExtractTransactionHandler
	Auth                       *AuthHandler
	Portfolio                  *PortfolioHandler
}

// InitHandlers wires up all dependencies and returns a Handlers struct
func InitHandlers(db *gorm.DB, cfg *config.Config) *Handlers {
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)

	// Initialize Price Service Manager
	priceServiceManager := provider.NewPriceServiceManager(cfg)

	// Initialize Portfolio Service
	portfolioService := services.NewPortfolioService(transactionRepo, priceServiceManager)

	// Initialize AI client once for reuse
	aiClient, err := ai.NewClient(cfg)
	if err != nil {
		panic("Failed to initialize AI client: " + err.Error())
	}

	return &Handlers{
		Transactions:               NewTransactionsHandler(transactionService),
		ExtractTransactionsHandler: NewExtractTransactionsHandler(cfg, aiClient),
		Auth:                       NewAuthHandler(db, cfg),
		Portfolio:                  NewPortfolioHandler(portfolioService),
	}
}
