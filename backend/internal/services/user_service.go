package services

import (
	"errors"
	"fmt"
	"time"

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

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) error {
	if err := s.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(user *models.User) error {
	if err := s.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// UpdateUserLastLogin updates the user's last login time
func (s *UserService) UpdateUserLastLogin(userID uint) error {
	now := time.Now()
	if err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", &now).Error; err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// DeactivateUser deactivates a user (soft delete)
func (s *UserService) DeactivateUser(userID uint) error {
	if err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	return nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(userID uint) error {
	if err := s.db.Delete(&models.User{}, userID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetActiveUsers retrieves all active users
func (s *UserService) GetActiveUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	query := s.db.Where("is_active = ? AND deleted_at IS NULL", true)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	return users, nil
}

// CountActiveUsers returns the count of active users
func (s *UserService) CountActiveUsers() (int64, error) {
	var count int64
	if err := s.db.Model(&models.User{}).Where("is_active = ? AND deleted_at IS NULL", true).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count active users: %w", err)
	}
	return count, nil
}
