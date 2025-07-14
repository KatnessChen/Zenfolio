package models

import (
	"time"

	"github.com/transaction-tracker/backend/internal/utils"
)

// JWTToken represents a JWT token record in the database
type JWTToken struct {
	ID         utils.UUID `gorm:"type:binary(16);primaryKey" json:"id"`
	UserID     utils.UUID `gorm:"type:binary(16);not null;index:idx_user_tokens" json:"user_id"`
	TokenHash  string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_token_hash" json:"token_hash"`
	IssuedAt   time.Time  `gorm:"not null" json:"issued_at"`
	ExpiresAt  time.Time  `gorm:"not null;index:idx_expires_at" json:"expires_at"`
	RevokedAt  *time.Time `gorm:"null" json:"revoked_at,omitempty"`
	LastUsedAt *time.Time `gorm:"null" json:"last_used_at,omitempty"`
	DeviceInfo string     `gorm:"type:json;null" json:"device_info,omitempty"`
	BaseModel

	// User relationship - foreign key is UserID pointing to users.user_id
	User User `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName returns the table name for JWTToken
func (JWTToken) TableName() string {
	return "jwt_tokens"
}

// IsExpired checks if the token is expired
func (j *JWTToken) IsExpired() bool {
	return time.Now().After(j.ExpiresAt)
}

// IsRevoked checks if the token is revoked
func (j *JWTToken) IsRevoked() bool {
	return j.RevokedAt != nil
}

// Revoke marks the token as revoked
func (j *JWTToken) Revoke() {
	now := time.Now()
	j.RevokedAt = &now
}

// UpdateLastUsed updates the last used timestamp
func (j *JWTToken) UpdateLastUsed() {
	now := time.Now()
	j.LastUsedAt = &now
}

// IsDeleted checks if the token is soft deleted
func (j *JWTToken) IsDeleted() bool {
	return j.DeletedAt.Valid
}
