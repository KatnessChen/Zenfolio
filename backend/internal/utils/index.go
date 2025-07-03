package utils

import (
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
