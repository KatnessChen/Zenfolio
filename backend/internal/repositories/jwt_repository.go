package repositories

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/transaction-tracker/backend/internal/models"
	"gorm.io/gorm"
)

// JWTRepository defines the interface for JWT token operations
type JWTRepository interface {
	Create(userID uint, tokenHash string, expiresAt time.Time, deviceInfo string) (*models.JWTToken, error)
	FindByTokenHash(tokenHash string) (*models.JWTToken, error)
	FindActiveTokensByUserID(userID uint) ([]models.JWTToken, error)
	UpdateLastUsed(tokenID string) error
	RevokeToken(tokenID string) error
	RevokeAllUserTokens(userID uint) error
	CleanupExpiredTokens() error
	GetUserTokenCount(userID uint) (int64, error)
}

// jwtRepository implements JWTRepository
type jwtRepository struct {
	db *gorm.DB
}

// NewJWTRepository creates a new JWT repository instance
func NewJWTRepository(db *gorm.DB) JWTRepository {
	return &jwtRepository{db: db}
}

// Create creates a new JWT token record
func (r *jwtRepository) Create(userID uint, tokenHash string, expiresAt time.Time, deviceInfo string) (*models.JWTToken, error) {
	token := &models.JWTToken{
		ID:         uuid.New().String(),
		UserID:     userID,
		TokenHash:  tokenHash,
		IssuedAt:   time.Now(),
		ExpiresAt:  expiresAt,
		DeviceInfo: deviceInfo,
	}

	if err := r.db.Create(token).Error; err != nil {
		return nil, fmt.Errorf("failed to create JWT token: %w", err)
	}

	return token, nil
}

// FindByTokenHash finds a JWT token by its hash
func (r *jwtRepository) FindByTokenHash(tokenHash string) (*models.JWTToken, error) {
	var token models.JWTToken

	err := r.db.Preload("User").Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to find token: %w", err)
	}

	return &token, nil
}

// FindActiveTokensByUserID finds all active tokens for a user
func (r *jwtRepository) FindActiveTokensByUserID(userID uint) ([]models.JWTToken, error) {
	var tokens []models.JWTToken

	err := r.db.Where("user_id = ? AND expires_at > ? AND revoked_at IS NULL",
		userID, time.Now()).Find(&tokens).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find active tokens: %w", err)
	}

	return tokens, nil
}

// UpdateLastUsed updates the last used timestamp for a token
func (r *jwtRepository) UpdateLastUsed(tokenID string) error {
	now := time.Now()
	err := r.db.Model(&models.JWTToken{}).
		Where("id = ?", tokenID).
		Update("last_used_at", now).Error

	if err != nil {
		return fmt.Errorf("failed to update last used: %w", err)
	}

	return nil
}

// RevokeToken revokes a specific token
func (r *jwtRepository) RevokeToken(tokenID string) error {
	now := time.Now()
	err := r.db.Model(&models.JWTToken{}).
		Where("id = ?", tokenID).
		Update("revoked_at", now).Error

	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (r *jwtRepository) RevokeAllUserTokens(userID uint) error {
	now := time.Now()
	err := r.db.Model(&models.JWTToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error

	if err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}

	return nil
}

// CleanupExpiredTokens removes expired tokens from the database
func (r *jwtRepository) CleanupExpiredTokens() error {
	err := r.db.Where("expires_at < ?", time.Now()).Delete(&models.JWTToken{}).Error
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return nil
}

// GetUserTokenCount returns the number of active tokens for a user
func (r *jwtRepository) GetUserTokenCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.JWTToken{}).
		Where("user_id = ? AND expires_at > ? AND revoked_at IS NULL",
			userID, time.Now()).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count user tokens: %w", err)
	}

	return count, nil
}

// HashToken creates a SHA-256 hash of the token for secure storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}
