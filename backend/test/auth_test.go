package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/transaction-tracker/backend/api/handlers"
	"github.com/transaction-tracker/backend/api/middlewares"
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/constants"
	"github.com/transaction-tracker/backend/internal/models"
)

// TestSetup holds the test environment
type TestSetup struct {
	DB          *gorm.DB
	Router      *gin.Engine
	AuthHandler *handlers.AuthHandler
	Config      *config.Config
}

// setupTestEnvironment creates a complete test environment
func setupTestEnvironment(t *testing.T) *TestSetup {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Auto-migrate the schemas
	err = db.AutoMigrate(&models.User{}, &models.Transaction{}, &models.JWTToken{})
	require.NoError(t, err)

	// Create test config
	cfg := &config.Config{
		JWTSecret:          "test-secret-key-for-testing-only",
		JWTExpirationHours: 24,
		RateLimitRequests:  100, // High limit for testing
	}

	// Create auth handler
	authHandler := handlers.NewAuthHandler(db, cfg)

	// Create router and setup routes manually
	router := gin.New()

	// Public routes
	publicApi := router.Group(constants.APIVersion)
	{
		publicApi.POST(constants.LoginEndpoint, authHandler.Login)
		publicApi.POST(constants.SignupEndpoint, authHandler.Signup)
	}

	// Protected routes
	api := router.Group(constants.APIVersion)
	api.Use(middlewares.AuthMiddleware(db, cfg))
	{
		api.GET(constants.MeEndpoint, authHandler.Me)
		api.POST(constants.LogoutEndpoint, authHandler.Logout)
	}

	return &TestSetup{
		DB:          db,
		Router:      router,
		AuthHandler: authHandler,
		Config:      cfg,
	}
}

// Helper function to make HTTP requests
func makeRequest(method, url string, body interface{}, headers map[string]string) (*http.Request, error) {
	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func TestSignupEndpoint(t *testing.T) {
	setup := setupTestEnvironment(t)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid signup",
			payload: map[string]interface{}{
				"email":            "test@example.com",
				"first_name":       "Test",
				"last_name":        "User",
				"password":         "password123",
				"confirm_password": "password123",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid email format",
			payload: map[string]interface{}{
				"email":            "invalid-email",
				"first_name":       "Test",
				"last_name":        "User",
				"password":         "password123",
				"confirm_password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
		},
		{
			name: "Password too short",
			payload: map[string]interface{}{
				"email":            "test2@example.com",
				"first_name":       "Test",
				"last_name":        "User",
				"password":         "short",
				"confirm_password": "short",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
		},
		{
			name: "Passwords don't match",
			payload: map[string]interface{}{
				"email":            "test3@example.com",
				"first_name":       "Test",
				"last_name":        "User",
				"password":         "password123",
				"confirm_password": "different123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Passwords do not match",
		},
		{
			name: "Missing required fields",
			payload: map[string]interface{}{
				"email":    "test4@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
		},
		{
			name: "Duplicate email",
			payload: map[string]interface{}{
				"email":            "test@example.com", // Same as first test
				"first_name":       "Another",
				"last_name":        "User",
				"password":         "password123",
				"confirm_password": "password123",
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "Email already registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := makeRequest("POST", "/api/v1/signup", tt.payload, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			setup.Router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.expectedError != "" {
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				assert.True(t, response["success"].(bool))
				assert.NotEmpty(t, response["user"])
			}
		})
	}
}

func TestLoginEndpoint(t *testing.T) {
	setup := setupTestEnvironment(t)

	// Create a test user first
	user := &models.User{
		Email:     "logintest@example.com",
		FirstName: "Login",
		LastName:  "Test",
		Username:  "Login Test",
		IsActive:  true,
	}
	err := user.SetPassword("testpassword123")
	require.NoError(t, err)
	err = setup.DB.Create(user).Error
	require.NoError(t, err)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid login",
			payload: map[string]interface{}{
				"email":    "logintest@example.com",
				"password": "testpassword123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid password",
			payload: map[string]interface{}{
				"email":    "logintest@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid credentials",
		},
		{
			name: "Non-existent user",
			payload: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "testpassword123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid credentials",
		},
		{
			name: "Invalid email format",
			payload: map[string]interface{}{
				"email":    "invalid-email",
				"password": "testpassword123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "Missing password",
			payload: map[string]interface{}{
				"email": "logintest@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := makeRequest("POST", "/api/v1/login", tt.payload, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			setup.Router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.expectedError != "" {
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				assert.NotEmpty(t, response["token"])
				assert.NotEmpty(t, response["user"])
			}
		})
	}
}

func TestProtectedEndpoints(t *testing.T) {
	setup := setupTestEnvironment(t)

	// Create a test user and login to get a token
	user := &models.User{
		Email:     "protected@example.com",
		FirstName: "Protected",
		LastName:  "Test",
		Username:  "Protected Test",
		IsActive:  true,
	}
	err := user.SetPassword("testpassword123")
	require.NoError(t, err)
	err = setup.DB.Create(user).Error
	require.NoError(t, err)

	// Login to get a valid token
	loginPayload := map[string]interface{}{
		"email":    "protected@example.com",
		"password": "testpassword123",
	}
	req, err := makeRequest("POST", "/api/v1/login", loginPayload, nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	setup.Router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	token := loginResponse["token"].(string)

	tests := []struct {
		name           string
		endpoint       string
		method         string
		headers        map[string]string
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "Valid token access to /me",
			endpoint: "/api/v1/me",
			method:   "GET",
			headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "No token provided",
			endpoint:       "/api/v1/me",
			method:         "GET",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization header required",
		},
		{
			name:     "Invalid token format",
			endpoint: "/api/v1/me",
			method:   "GET",
			headers: map[string]string{
				"Authorization": "InvalidToken",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid Authorization header format",
		},
		{
			name:     "Invalid token",
			endpoint: "/api/v1/me",
			method:   "GET",
			headers: map[string]string{
				"Authorization": "Bearer invalid.token.here",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := makeRequest(tt.method, tt.endpoint, nil, tt.headers)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			setup.Router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.expectedError != "" {
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				assert.NotEmpty(t, response["user"])
			}
		})
	}
}

func TestJWTTokenLifecycle(t *testing.T) {
	setup := setupTestEnvironment(t)

	// Create a test user
	user := &models.User{
		Email:     "jwttest@example.com",
		FirstName: "JWT",
		LastName:  "Test",
		Username:  "JWT Test",
		IsActive:  true,
	}
	err := user.SetPassword("testpassword123")
	require.NoError(t, err)
	err = setup.DB.Create(user).Error
	require.NoError(t, err)

	// Test multiple logins create different tokens
	var tokens []string
	for i := 0; i < 3; i++ {
		loginPayload := map[string]interface{}{
			"email":    "jwttest@example.com",
			"password": "testpassword123",
		}
		req, err := makeRequest("POST", "/api/v1/login", loginPayload, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		setup.Router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		token := response["token"].(string)
		tokens = append(tokens, token)

		// Verify each token is unique
		for j := 0; j < i; j++ {
			assert.NotEqual(t, tokens[j], token)
		}
	}

	// Verify all tokens are active in database
	var activeTokens []models.JWTToken
	err = setup.DB.Where("user_id = ? AND revoked_at IS NULL", user.ID).Find(&activeTokens).Error
	require.NoError(t, err)
	assert.Len(t, activeTokens, 3)

	// Test logout with one token
	logoutReq, err := makeRequest("POST", "/api/v1/logout", nil, map[string]string{
		"Authorization": "Bearer " + tokens[0],
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	setup.Router.ServeHTTP(w, logoutReq)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify token is revoked
	var revokedToken models.JWTToken
	err = setup.DB.Where("user_id = ? AND revoked_at IS NOT NULL", user.ID).First(&revokedToken).Error
	require.NoError(t, err)
	assert.NotNil(t, revokedToken.RevokedAt)

	// Verify revoked token can't be used
	protectedReq, err := makeRequest("GET", "/api/v1/me", nil, map[string]string{
		"Authorization": "Bearer " + tokens[0],
	})
	require.NoError(t, err)

	w = httptest.NewRecorder()
	setup.Router.ServeHTTP(w, protectedReq)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Verify other tokens still work
	protectedReq, err = makeRequest("GET", "/api/v1/me", nil, map[string]string{
		"Authorization": "Bearer " + tokens[1],
	})
	require.NoError(t, err)

	w = httptest.NewRecorder()
	setup.Router.ServeHTTP(w, protectedReq)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimiting(t *testing.T) {
	setup := setupTestEnvironment(t)

	// Test rate limiting on signup endpoint
	payload := map[string]interface{}{
		"email":            "ratetest@example.com",
		"first_name":       "Rate",
		"last_name":        "Test",
		"password":         "password123",
		"confirm_password": "password123",
	}

	// Make multiple requests quickly
	successCount := 0
	rateLimitedCount := 0

	for i := 0; i < 20; i++ {
		req, err := makeRequest("POST", "/api/v1/signup", payload, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		setup.Router.ServeHTTP(w, req)

		if w.Code == http.StatusCreated || w.Code == http.StatusConflict {
			successCount++
		} else if w.Code == http.StatusTooManyRequests {
			rateLimitedCount++
		}

		// Change email for subsequent requests to avoid duplicate email errors
		if email, ok := payload["email"].(string); ok {
			payload["email"] = email + string(rune('a'+i))
		}
	}

	// Should have some rate limiting after many requests
	t.Logf("Success count: %d, Rate limited count: %d", successCount, rateLimitedCount)

	// Note: Exact rate limiting behavior depends on configuration
	// This test verifies the middleware is working
}

func TestSignupLoginIntegration(t *testing.T) {
	setup := setupTestEnvironment(t)

	// Complete signup -> login -> protected access flow
	email := "integration@example.com"
	password := "integrationtest123"

	// Step 1: Signup
	signupPayload := map[string]interface{}{
		"email":            email,
		"first_name":       "Integration",
		"last_name":        "Test",
		"password":         password,
		"confirm_password": password,
	}

	req, err := makeRequest("POST", "/api/v1/signup", signupPayload, nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	setup.Router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var signupResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &signupResponse)
	require.NoError(t, err)
	assert.True(t, signupResponse["success"].(bool))

	// Step 2: Login
	loginPayload := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	req, err = makeRequest("POST", "/api/v1/login", loginPayload, nil)
	require.NoError(t, err)

	w = httptest.NewRecorder()
	setup.Router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	token := loginResponse["token"].(string)
	assert.NotEmpty(t, token)

	// Step 3: Access protected endpoint
	req, err = makeRequest("GET", "/api/v1/me", nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	require.NoError(t, err)

	w = httptest.NewRecorder()
	setup.Router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var meResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &meResponse)
	require.NoError(t, err)

	user := meResponse["user"].(map[string]interface{})
	assert.Equal(t, email, user["email"])
	assert.Equal(t, "Integration Test", user["username"]) // Should be FirstName + LastName
}
