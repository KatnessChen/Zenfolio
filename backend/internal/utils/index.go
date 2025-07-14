package utils

import (
	"github.com/google/uuid"
	"github.com/transaction-tracker/backend/internal/types"
)

// TradeTypeFromString converts string to TradeType
func TradeTypeFromString(s string) (types.TradeType, bool) {
	switch s {
	case "Buy":
		return types.TradeTypeBuy, true
	case "Sell":
		return types.TradeTypeSell, true
	case "Dividends":
		return types.TradeTypeDividend, true
	default:
		return "", false
	}
}

// GenerateUUID generates a new UUID string (v4)
func GenerateUUID() string {
	return uuid.New().String()
}
