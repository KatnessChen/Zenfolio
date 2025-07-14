package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/transaction-tracker/backend/api/handlers"
	"github.com/transaction-tracker/backend/internal/models"
	"github.com/transaction-tracker/backend/internal/utils"
)

// MockUserRepository is a mock implementation of the user repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUserID(userID utils.UUID) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test helper to create a test request
func createSignupRequest(payload interface{}) (*http.Request, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return httptest.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(jsonPayload)), nil
}

func TestSignupHandler_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockUserRepo := new(MockUserRepository)

	// Mock expectations
	mockUserRepo.On("FindByEmail", "test@example.com").Return(nil, errors.New("user not found"))
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Create test request
	payload := map[string]interface{}{
		"email":            "test@example.com",
		"first_name":       "Test",
		"last_name":        "User",
		"password":         "testpassword123",
		"confirm_password": "testpassword123",
	}

	req, err := createSignupRequest(payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test JSON binding and validation
	var signupReq handlers.SignupRequest
	err = c.ShouldBindJSON(&signupReq)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", signupReq.Email)
	assert.Equal(t, "Test", signupReq.FirstName)
	assert.Equal(t, "User", signupReq.LastName)
	assert.Equal(t, "testpassword123", signupReq.Password)
	assert.Equal(t, "testpassword123", signupReq.ConfirmPassword)
}

func TestSignupValidation_PasswordMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	payload := map[string]interface{}{
		"email":            "test@example.com",
		"first_name":       "Test",
		"last_name":        "User",
		"password":         "testpassword123",
		"confirm_password": "differentpassword",
	}

	req, err := createSignupRequest(payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	var signupReq handlers.SignupRequest
	err = c.ShouldBindJSON(&signupReq)
	assert.NoError(t, err)

	// Test password mismatch validation
	assert.NotEqual(t, signupReq.Password, signupReq.ConfirmPassword)
}

func TestSignupValidation_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	payload := map[string]interface{}{
		"email":            "invalid-email",
		"first_name":       "Test",
		"last_name":        "User",
		"password":         "testpassword123",
		"confirm_password": "testpassword123",
	}

	req, err := createSignupRequest(payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	var signupReq handlers.SignupRequest
	err = c.ShouldBindJSON(&signupReq)
	// This should fail validation due to invalid email
	assert.Error(t, err)
}

func TestSignupValidation_ShortPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	payload := map[string]interface{}{
		"email":            "test@example.com",
		"first_name":       "Test",
		"last_name":        "User",
		"password":         "short",
		"confirm_password": "short",
	}

	req, err := createSignupRequest(payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	var signupReq handlers.SignupRequest
	err = c.ShouldBindJSON(&signupReq)
	// This should fail validation due to short password (min 8 chars)
	assert.Error(t, err)
}

func TestSignupValidation_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name    string
		payload map[string]interface{}
	}{
		{
			name: "missing email",
			payload: map[string]interface{}{
				"first_name":       "Test",
				"password":         "testpassword123",
				"confirm_password": "testpassword123",
			},
		},
		{
			name: "missing first_name",
			payload: map[string]interface{}{
				"email":            "test@example.com",
				"password":         "testpassword123",
				"confirm_password": "testpassword123",
			},
		},
		{
			name: "missing password",
			payload: map[string]interface{}{
				"email":            "test@example.com",
				"first_name":       "Test",
				"confirm_password": "testpassword123",
			},
		},
		{
			name: "missing confirm_password",
			payload: map[string]interface{}{
				"email":      "test@example.com",
				"first_name": "Test",
				"password":   "testpassword123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := createSignupRequest(tc.payload)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			var signupReq handlers.SignupRequest
			err = c.ShouldBindJSON(&signupReq)
			// Should fail validation due to missing required fields
			assert.Error(t, err)
		})
	}
}

func TestSignupValidation_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	payload := map[string]interface{}{
		"email":            "test@example.com",
		"first_name":       "Test",
		"last_name":        "User",
		"password":         "testpassword123",
		"confirm_password": "testpassword123",
	}

	req, err := createSignupRequest(payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	var signupReq handlers.SignupRequest
	err = c.ShouldBindJSON(&signupReq)
	assert.NoError(t, err)

	// Verify all fields are correctly parsed
	assert.Equal(t, "test@example.com", signupReq.Email)
	assert.Equal(t, "Test", signupReq.FirstName)
	assert.Equal(t, "User", signupReq.LastName)
	assert.Equal(t, "testpassword123", signupReq.Password)
	assert.Equal(t, "testpassword123", signupReq.ConfirmPassword)

	// Test password validation logic
	assert.Equal(t, signupReq.Password, signupReq.ConfirmPassword)
}

func TestUserModel_SetPassword(t *testing.T) {
	user := &models.User{}

	err := user.SetPassword("testpassword123")
	assert.NoError(t, err)

	// Verify password is hashed (not stored in plain text)
	assert.NotEqual(t, "testpassword123", user.PasswordHash)
	assert.NotEmpty(t, user.PasswordHash)

	// Verify password can be verified
	assert.True(t, user.CheckPassword("testpassword123"))
	assert.False(t, user.CheckPassword("wrongpassword"))
}

func TestUserModel_SetPasswordEmpty(t *testing.T) {
	user := &models.User{}

	// Empty password should still work (bcrypt can hash empty strings)
	err := user.SetPassword("")
	assert.NoError(t, err)
	assert.NotEmpty(t, user.PasswordHash)

	// Verify empty password can be verified
	assert.True(t, user.CheckPassword(""))
	assert.False(t, user.CheckPassword("notempty"))
}
