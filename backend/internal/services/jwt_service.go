package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/models"
	"github.com/transaction-tracker/backend/internal/repositories"
)

// JWTClaims represents the claims in our JWT token
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	TokenID  string `json:"token_id"`
	jwt.RegisteredClaims
}

// DeviceInfo represents device information for security tracking
type DeviceInfo struct {
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
	Browser   string `json:"browser,omitempty"`
	OS        string `json:"os,omitempty"`
}

// JWTService defines the interface for JWT operations
type JWTService interface {
	GenerateToken(user *models.User, deviceInfo DeviceInfo) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	RefreshToken(tokenString string, deviceInfo DeviceInfo) (string, error)
	RevokeToken(tokenString string) error
	RevokeAllUserTokens(userID uint) error
	GetActiveTokens(userID uint) ([]models.JWTToken, error)
	CleanupExpiredTokens() error
	ExtractTokenID(tokenString string) (string, error)
}

// jwtService implements JWTService
type jwtService struct {
	cfg           *config.Config
	jwtRepository repositories.JWTRepository
}

// NewJWTService creates a new JWT service instance
func NewJWTService(cfg *config.Config, jwtRepo repositories.JWTRepository) JWTService {
	return &jwtService{
		cfg:           cfg,
		jwtRepository: jwtRepo,
	}
}

// GenerateToken generates a new JWT token and stores it in the database
func (s *jwtService) GenerateToken(user *models.User, deviceInfo DeviceInfo) (string, error) {
	// Validate input
	if user == nil {
		return "", fmt.Errorf("user cannot be nil")
	}

	// Calculate expiration time (24 hours from now)
	expiresAt := time.Now().Add(time.Duration(s.cfg.JWTExpirationHours) * time.Hour)

	// Create device info JSON
	deviceInfoJSON, err := json.Marshal(deviceInfo)
	if err != nil {
		return "", fmt.Errorf("failed to marshal device info: %w", err)
	}

	// Generate a unique token ID using UUID
	tokenID := uuid.New().String()

	// Create JWT claims
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		TokenID:  tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "transaction-tracker",
			Subject:   fmt.Sprintf("user:%d", user.ID),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	// Create composite hash including both token string and token ID
	// This ensures the database record is linked to the specific token ID in the JWT
	compositeString := fmt.Sprintf("%s:%s", tokenString, tokenID)
	tokenHash := repositories.HashToken(compositeString)

	// Store token in database
	_, err = s.jwtRepository.Create(user.ID, tokenHash, expiresAt, string(deviceInfoJSON))
	if err != nil {
		return "", fmt.Errorf("failed to store token in database: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and checks against the database
func (s *jwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check token in database
	// Extract token ID from claims to create composite hash
	tokenID := claims.TokenID
	if tokenID == "" {
		return nil, fmt.Errorf("token ID not found in claims")
	}

	// Create composite hash including both token string and token ID
	compositeString := fmt.Sprintf("%s:%s", tokenString, tokenID)
	tokenHash := repositories.HashToken(compositeString)

	jwtToken, err := s.jwtRepository.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("token not found in database: %w", err)
	}

	// Validate token status - check if expired or revoked
	if jwtToken.IsExpired() {
		return nil, fmt.Errorf("token has expired")
	}

	if jwtToken.IsRevoked() {
		return nil, fmt.Errorf("token has been revoked")
	}

	// Update last used timestamp
	if err := s.jwtRepository.UpdateLastUsed(jwtToken.ID); err != nil {
		// Log error but don't fail validation
		fmt.Printf("Warning: failed to update last used timestamp: %v\n", err)
	}

	return claims, nil
}

// RefreshToken refreshes an existing token (generates a new one and revokes the old)
func (s *jwtService) RefreshToken(tokenString string, deviceInfo DeviceInfo) (string, error) {
	// Validate current token
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("cannot refresh invalid token: %w", err)
	}

	// Create composite hash to find the token in database
	compositeString := fmt.Sprintf("%s:%s", tokenString, claims.TokenID)
	tokenHash := repositories.HashToken(compositeString)

	jwtToken, err := s.jwtRepository.FindByTokenHash(tokenHash)
	if err != nil {
		return "", fmt.Errorf("token not found: %w", err)
	}

	// Revoke current token
	if err := s.jwtRepository.RevokeToken(jwtToken.ID); err != nil {
		return "", fmt.Errorf("failed to revoke current token: %w", err)
	}

	// Generate new token
	return s.GenerateToken(&jwtToken.User, deviceInfo)
}

// RevokeToken revokes a specific token
func (s *jwtService) RevokeToken(tokenString string) error {
	// Extract token ID from the JWT claims
	tokenID, err := s.ExtractTokenID(tokenString)
	if err != nil {
		return fmt.Errorf("failed to extract token ID: %w", err)
	}

	// Create composite hash to find the token in database
	compositeString := fmt.Sprintf("%s:%s", tokenString, tokenID)
	tokenHash := repositories.HashToken(compositeString)

	jwtToken, err := s.jwtRepository.FindByTokenHash(tokenHash)
	if err != nil {
		return fmt.Errorf("token not found: %w", err)
	}

	return s.jwtRepository.RevokeToken(jwtToken.ID)
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (s *jwtService) RevokeAllUserTokens(userID uint) error {
	return s.jwtRepository.RevokeAllUserTokens(userID)
}

// GetActiveTokens returns all active tokens for a user
func (s *jwtService) GetActiveTokens(userID uint) ([]models.JWTToken, error) {
	return s.jwtRepository.FindActiveTokensByUserID(userID)
}

// CleanupExpiredTokens removes expired tokens from the database
func (s *jwtService) CleanupExpiredTokens() error {
	return s.jwtRepository.CleanupExpiredTokens()
}

// ExtractTokenID extracts the token ID from a JWT string without full validation
// Useful for logging, tracking, and audit purposes
func (s *jwtService) ExtractTokenID(tokenString string) (string, error) {
	// Parse token without signature validation to extract claims
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && claims.TokenID != "" {
		return claims.TokenID, nil
	}

	return "", fmt.Errorf("token ID not found in claims")
}
