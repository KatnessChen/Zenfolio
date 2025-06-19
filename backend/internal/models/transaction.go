package models

import (
	"time"

	"github.com/transaction-tracker/backend/internal/types"
	"gorm.io/gorm"
)

// Transaction represents a financial transaction
type Transaction struct {
	BaseModel
	UserID          uint            `gorm:"not null;index" json:"user_id"`
	Type            types.TradeType `gorm:"size:50;not null;index" json:"type"`
	Symbol          string          `gorm:"size:20;not null;index" json:"symbol"`
	Quantity        float64         `gorm:"type:decimal(15,4);not null" json:"quantity"`
	Price           float64         `gorm:"type:decimal(15,4);not null" json:"price"`
	Amount          float64         `gorm:"type:decimal(15,2);not null" json:"amount"`
	Currency        string          `gorm:"size:3;not null;default:'USD'" json:"currency"`
	Broker          string          `gorm:"size:100" json:"broker"`
	Account         string          `gorm:"size:100" json:"account"`
	TransactionDate time.Time       `gorm:"not null;index" json:"transaction_date"`
	UserNotes       string          `gorm:"type:text" json:"user_notes"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user,omitempty"`
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate hook for Transaction model
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate hook for Transaction model
func (t *Transaction) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
