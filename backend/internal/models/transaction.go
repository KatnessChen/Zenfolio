package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/transaction-tracker/backend/internal/types"
	"github.com/transaction-tracker/backend/internal/utils"
	"gorm.io/gorm"
)

// Transaction represents a financial transaction
type Transaction struct {
	TransactionID   utils.UUID      `gorm:"type:binary(16);primaryKey" json:"transaction_id"`
	UserID          utils.UUID      `gorm:"type:binary(16);not null;index" json:"user_id"`
	TradeType       types.TradeType `gorm:"column:trade_type;size:50;not null;index" json:"trade_type"`
	Symbol          string          `gorm:"size:20;not null;index" json:"symbol"`
	Quantity        float64         `gorm:"type:decimal(15,4);not null" json:"quantity"`
	Price           float64         `gorm:"type:decimal(15,4);not null" json:"price"`
	Amount          float64         `gorm:"type:decimal(15,2);not null" json:"amount"`
	Currency        string          `gorm:"size:3;not null;default:'USD'" json:"currency"`
	Exchange        string          `gorm:"size:50" json:"exchange"`
	Broker          string          `gorm:"size:100" json:"broker"`
	TransactionDate time.Time       `gorm:"not null;index" json:"transaction_date"`
	UserNotes       string          `gorm:"type:text" json:"user_notes"`
	BaseModel

	// User relationship - foreign key is UserID pointing to users.user_id
	User User `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate hook for Transaction model
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.TransactionID.UUID == uuid.Nil {
		t.TransactionID = utils.UUID{UUID: uuid.New()}
	}
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
