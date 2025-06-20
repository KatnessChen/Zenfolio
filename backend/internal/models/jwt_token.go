package models

import (
	"time"
)

// JWTToken represents a JWT token record in the database
type JWTToken struct {
	BaseModel
	ID         string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index:idx_user_tokens" json:"user_id"`
	TokenHash  string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"token_hash"`
	IssuedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"issued_at"`
	ExpiresAt  time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt  *time.Time `gorm:"null" json:"revoked_at,omitempty"`
	LastUsedAt *time.Time `gorm:"null" json:"last_used_at,omitempty"`
	DeviceInfo string     `gorm:"type:json;null" json:"device_info,omitempty"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
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
