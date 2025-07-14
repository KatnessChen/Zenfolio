package repositories

import (
	"fmt"

	"github.com/transaction-tracker/backend/internal/models"
	"github.com/transaction-tracker/backend/internal/utils"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user operations
type UserRepository interface {
	Create(user *models.User) error
	FindByUserID(userID utils.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// FindByUserID finds a user by user ID (UUID)
func (r *userRepository) FindByUserID(userID utils.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
