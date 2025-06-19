package services

import (
	"fmt"

	"github.com/transaction-tracker/backend/internal/models"
	"gorm.io/gorm"
)

// UserService handles user-related database operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CountActiveUsers returns the count of active users
func (s *UserService) CountActiveUsers() (int64, error) {
	var count int64
	if err := s.db.Model(&models.User{}).Where("is_active = ? AND deleted_at IS NULL", true).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count active users: %w", err)
	}
	return count, nil
}
