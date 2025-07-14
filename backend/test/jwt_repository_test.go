package test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testutils "github.com/transaction-tracker/backend/testutils"
	"gorm.io/gorm"

	"github.com/transaction-tracker/backend/internal/models"
	"github.com/transaction-tracker/backend/internal/repositories"
)

// setupJWTRepositoryTest sets up the test environment for JWT repository tests
func setupJWTRepositoryTest(t *testing.T) (*gorm.DB, repositories.JWTRepository, *models.User) {
	// Use shared MySQL test DB
	db := testutils.SetupTestDB(t)

	jwtRepo := repositories.NewJWTRepository(db)

	// Create a test user
	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	err := testUser.SetPassword("password123")
	require.NoError(t, err)

	err = db.Create(testUser).Error
	require.NoError(t, err)

	return db, jwtRepo, testUser
}

// Test 2.1: Token Storage Tests
func TestJWTRepository_Create(t *testing.T) {
	db, jwtRepo, testUser := setupJWTRepositoryTest(t)

	tokenHash := "test_hash_" + uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	deviceInfo := `{"device": "test", "browser": "test"}`

	token, err := jwtRepo.Create(testUser.UserID, tokenHash, expiresAt, deviceInfo)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.ID)
	assert.Equal(t, testUser.UserID, token.UserID)
	assert.Equal(t, tokenHash, token.TokenHash)

	// Verify token was stored in database
	var storedToken models.JWTToken
	err = db.Where("id = ?", token.ID).First(&storedToken).Error
	assert.NoError(t, err)
	assert.Equal(t, token.UserID, storedToken.UserID)
	assert.Equal(t, token.TokenHash, storedToken.TokenHash)
}

func TestJWTRepository_Create_DuplicateTokenHash(t *testing.T) {
	_, jwtRepo, testUser := setupJWTRepositoryTest(t)

	tokenHash := "duplicate_hash"
	expiresAt := time.Now().Add(24 * time.Hour)
	deviceInfo := `{"device": "test", "browser": "test"}`

	// Create first token
	_, err := jwtRepo.Create(testUser.UserID, tokenHash, expiresAt, deviceInfo)
	assert.NoError(t, err)

	// Try to create second token with same hash
	_, err = jwtRepo.Create(testUser.UserID, tokenHash, expiresAt, deviceInfo)
	assert.Error(t, err) // Should fail due to unique constraint
}

// Test 2.2: Token Retrieval Tests
func TestJWTRepository_FindByTokenHash(t *testing.T) {
	_, jwtRepo, testUser := setupJWTRepositoryTest(t)

	tokenHash := "test_hash_" + uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	deviceInfo := `{"device": "test", "browser": "test"}`

	token, err := jwtRepo.Create(testUser.UserID, tokenHash, expiresAt, deviceInfo)
	require.NoError(t, err)

	// Test successful retrieval
	retrievedToken, err := jwtRepo.FindByTokenHash(token.TokenHash)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedToken)
	assert.Equal(t, token.ID, retrievedToken.ID)
	assert.Equal(t, token.UserID, retrievedToken.UserID)
}

func TestJWTRepository_FindByTokenHash_NotFound(t *testing.T) {
	_, jwtRepo, _ := setupJWTRepositoryTest(t)

	// Test with non-existent hash
	retrievedToken, err := jwtRepo.FindByTokenHash("non_existent_hash")
	assert.Error(t, err)
	assert.Nil(t, retrievedToken)
	assert.Contains(t, err.Error(), "token not found")
}

func TestJWTRepository_FindByTokenHash_RevokedToken(t *testing.T) {
	db, jwtRepo, testUser := setupJWTRepositoryTest(t)

	tokenHash := "test_hash_" + uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	deviceInfo := `{"device": "test", "browser": "test"}`

	token, err := jwtRepo.Create(testUser.UserID, tokenHash, expiresAt, deviceInfo)
	require.NoError(t, err)

	// Manually revoke the token
	token.Revoke()
	err = db.Save(token).Error
	require.NoError(t, err)

	// The repository still retrieves revoked tokens, but the token will have RevokedAt set
	retrievedToken, err := jwtRepo.FindByTokenHash(token.TokenHash)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedToken)
	assert.True(t, retrievedToken.IsRevoked())
}

// Test 2.3: Token Revocation Tests
func TestJWTRepository_RevokeToken(t *testing.T) {
	db, jwtRepo, testUser := setupJWTRepositoryTest(t)

	tokenHash := "test_hash_" + uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	deviceInfo := `{"device": "test", "browser": "test"}`

	token, err := jwtRepo.Create(testUser.UserID, tokenHash, expiresAt, deviceInfo)
	require.NoError(t, err)

	// Revoke the token
	err = jwtRepo.RevokeToken(token.ID)
	assert.NoError(t, err)

	// Verify token is revoked in database
	var revokedToken models.JWTToken
	err = db.Where("id = ?", token.ID).First(&revokedToken).Error
	assert.NoError(t, err)
	assert.True(t, revokedToken.IsRevoked())
	assert.NotNil(t, revokedToken.RevokedAt)
}

func TestJWTRepository_RevokeToken_NotFound(t *testing.T) {
	_, jwtRepo, _ := setupJWTRepositoryTest(t)

	// Try to revoke non-existent token - the repository doesn't check if token exists
	// It just updates 0 rows, which is not an error in GORM
	nonExistentID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	err := jwtRepo.RevokeToken(nonExistentID)
	assert.NoError(t, err) // This succeeds but affects 0 rows
}

// Test 2.4: Token Listing Tests
func TestJWTRepository_FindActiveTokensByUserID(t *testing.T) {
	_, jwtRepo, testUser := setupJWTRepositoryTest(t)

	// Create multiple tokens: some active, some revoked
	for i := 0; i < 2; i++ {
		tokenHash := "active_hash_" + uuid.New().String()
		expiresAt := time.Now().Add(24 * time.Hour)
		deviceInfo := `{"device": "test", "browser": "test"}`

		_, err := jwtRepo.Create(testUser.UserID, tokenHash, expiresAt, deviceInfo)
		require.NoError(t, err)
	}

	// Create and revoke one token
	revokedHash := "revoked_hash_" + uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	deviceInfo := `{"device": "test", "browser": "test"}`

	revokedToken, err := jwtRepo.Create(testUser.UserID, revokedHash, expiresAt, deviceInfo)
	require.NoError(t, err)
	err = jwtRepo.RevokeToken(revokedToken.ID)
	require.NoError(t, err)

	// Get active tokens
	retrievedTokens, err := jwtRepo.FindActiveTokensByUserID(testUser.UserID)
	assert.NoError(t, err)
	assert.Len(t, retrievedTokens, 2)

	// Verify only active tokens are returned
	for _, token := range retrievedTokens {
		assert.False(t, token.IsRevoked())
		assert.Equal(t, testUser.UserID, token.UserID)
	}
}

func TestJWTRepository_FindActiveTokensByUserID_NoTokens(t *testing.T) {
	_, jwtRepo, testUser := setupJWTRepositoryTest(t)

	// Get active tokens for user with no tokens
	retrievedTokens, err := jwtRepo.FindActiveTokensByUserID(testUser.UserID)
	assert.NoError(t, err)
	assert.Empty(t, retrievedTokens)
}

// Test 2.5: Token Cleanup Tests
func TestJWTRepository_CleanupExpiredTokens(t *testing.T) {
	db, jwtRepo, testUser := setupJWTRepositoryTest(t)

	// Create expired token
	expiredHash := "expired_hash_" + uuid.New().String()
	expiredTime := time.Now().Add(-1 * time.Hour) // Expired 1 hour ago
	deviceInfo := `{"device": "test", "browser": "test"}`

	_, err := jwtRepo.Create(testUser.UserID, expiredHash, expiredTime, deviceInfo)
	require.NoError(t, err)

	// Create valid token
	validHash := "valid_hash_" + uuid.New().String()
	validTime := time.Now().Add(24 * time.Hour) // Expires in 24 hours

	validToken, err := jwtRepo.Create(testUser.UserID, validHash, validTime, deviceInfo)
	require.NoError(t, err)

	// Delete expired tokens
	err = jwtRepo.CleanupExpiredTokens()
	assert.NoError(t, err)

	// Verify expired token is deleted, valid token remains
	var tokens []models.JWTToken
	err = db.Find(&tokens).Error
	assert.NoError(t, err)
	assert.Len(t, tokens, 1)
	assert.Equal(t, validToken.ID, tokens[0].ID)
}

func TestJWTRepository_CleanupExpiredTokens_NoExpiredTokens(t *testing.T) {
	_, jwtRepo, testUser := setupJWTRepositoryTest(t)

	// Create only valid tokens
	for i := 0; i < 3; i++ {
		tokenHash := "valid_hash_" + uuid.New().String()
		expiresAt := time.Now().Add(24 * time.Hour)
		deviceInfo := `{"device": "test", "browser": "test"}`

		_, err := jwtRepo.Create(testUser.UserID, tokenHash, expiresAt, deviceInfo)
		require.NoError(t, err)
	}

	// Try to cleanup expired tokens
	err := jwtRepo.CleanupExpiredTokens()
	assert.NoError(t, err) // Should succeed even with no expired tokens
}

// Test 2.6: Maintenance Tests

func TestJWTRepository_UpdateLastUsed(t *testing.T) {
	db, jwtRepo, testUser := setupJWTRepositoryTest(t)

	tokenHash := "test_hash_" + uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	deviceInfo := `{"device": "test", "browser": "test"}`

	token, err := jwtRepo.Create(testUser.UserID, tokenHash, expiresAt, deviceInfo)
	require.NoError(t, err)

	// Initially no last used timestamp
	assert.Nil(t, token.LastUsedAt)

	// Update last used
	err = jwtRepo.UpdateLastUsed(token.ID)
	assert.NoError(t, err)

	// Verify last used timestamp is set
	var updatedToken models.JWTToken
	err = db.Where("id = ?", token.ID).First(&updatedToken).Error
	assert.NoError(t, err)
	assert.NotNil(t, updatedToken.LastUsedAt)
}

// Test 2.7: Edge Cases and Error Conditions
func TestJWTRepository_TokenWithLargeDeviceInfo(t *testing.T) {
	_, jwtRepo, testUser := setupJWTRepositoryTest(t)

	// Create token with large device info
	tokenHash := "test_hash_" + uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	// Create a large string with valid JSON characters
	largeString := strings.Repeat("X", 1000)
	largeDeviceInfo := `{"device": "` + largeString + `", "browser": "test"}`

	token, err := jwtRepo.Create(testUser.UserID, tokenHash, expiresAt, largeDeviceInfo)
	assert.NoError(t, err)

	// Retrieve and verify
	retrievedToken, err := jwtRepo.FindByTokenHash(token.TokenHash)
	assert.NoError(t, err)
	assert.Equal(t, largeDeviceInfo, retrievedToken.DeviceInfo)
}
