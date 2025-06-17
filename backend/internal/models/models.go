package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common columns for all tables
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User represents a user in the system
type User struct {
	BaseModel
	Username     string        `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email        string        `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string        `gorm:"size:255;not null" json:"-"`
	FirstName    string        `gorm:"size:100" json:"first_name"`
	LastName     string        `gorm:"size:100" json:"last_name"`
	IsActive     bool          `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time    `json:"last_login_at"`
	
	// Relationships
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
}

// Transaction represents a financial transaction
type Transaction struct {
	BaseModel
	UserID          uint                   `gorm:"not null;index" json:"user_id"`
	Type            string                 `gorm:"size:50;not null;index" json:"type"` // buy, sell, dividend, etc.
	Symbol          string                 `gorm:"size:20;not null;index" json:"symbol"`
	Quantity        float64                `gorm:"type:decimal(15,4);not null" json:"quantity"`
	Price           float64                `gorm:"type:decimal(15,4);not null" json:"price"`
	Amount          float64                `gorm:"type:decimal(15,2);not null" json:"amount"`
	Fee             float64                `gorm:"type:decimal(15,2);default:0" json:"fee"`
	Currency        string                 `gorm:"size:3;not null;default:'USD'" json:"currency"`
	Broker          string                 `gorm:"size:100" json:"broker"`
	Account         string                 `gorm:"size:100" json:"account"`
	TransactionDate time.Time              `gorm:"not null;index" json:"transaction_date"`
	SettlementDate  *time.Time             `json:"settlement_date"`
	Description     string                 `gorm:"type:text" json:"description"`
	Reference       string                 `gorm:"size:255" json:"reference"`
	Status          string                 `gorm:"size:20;default:'completed';index" json:"status"`
	Tags            string                 `gorm:"type:text" json:"tags"` // JSON array of tags
	Metadata        string                 `gorm:"type:json" json:"metadata"` // Additional metadata as JSON
	
	// AI Processing
	ExtractedFrom   string `gorm:"size:255" json:"extracted_from"` // Source of extraction (file, manual, etc.)
	ProcessingNotes string `gorm:"type:text" json:"processing_notes"`
	
	// Relationships
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate hook for User model
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate hook for User model
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
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
