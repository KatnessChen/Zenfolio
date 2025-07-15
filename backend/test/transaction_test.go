package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/transaction-tracker/backend/api/handlers"
	"github.com/transaction-tracker/backend/internal/models"
	"github.com/transaction-tracker/backend/internal/types"
	"gorm.io/gorm"
)

func createTestUserWithUsername(db *gorm.DB, username, email string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Email:    email,
	}
	// Set password using the model's method
	if err := user.SetPassword("test123"); err != nil {
		return nil, err
	}
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func TestDeleteTransaction(t *testing.T) {
	// Setup test handler and database
	transactionsHandler, db, err := setupTestTransactionsHandler(t)
	require.NoError(t, err)

	// Create test user
	user, err := createTestUser(db, "testuser@example.com")
	require.NoError(t, err)

	// Create a test transaction
	transaction := &models.Transaction{
		UserID:          user.UserID,
		Symbol:          "AAPL",
		TradeType:       types.TradeType("Buy"),
		Quantity:        100,
		Price:           150.00,
		Amount:          15000.00,
		Currency:        "USD",
		Broker:          "Test Broker",
		Exchange:        "NASDAQ",
		TransactionDate: time.Now(),
		UserNotes:       "Test transaction for deletion",
	}
	err = db.Create(transaction).Error
	require.NoError(t, err)

	// Setup gin context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/transaction-history/:id", func(c *gin.Context) {
		c.Set("user_id", user.UserID)
		transactionsHandler.DeleteTransaction(c)
	})

	// Test successful deletion
	t.Run("Success_DeleteTransaction", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/transaction-history/%s", transaction.TransactionID.String()), nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.DeleteTransactionResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Equal(t, "Transaction deleted successfully", response.Message)
		assert.Equal(t, []string{transaction.TransactionID.String()}, response.Data.DeletedIDs)

		// Verify transaction is soft deleted
		var deletedTransaction models.Transaction
		err = db.Unscoped().Where("transaction_id = ?", transaction.TransactionID).First(&deletedTransaction).Error
		require.NoError(t, err)
		assert.NotNil(t, deletedTransaction.DeletedAt)
	})

	// Test transaction not found
	t.Run("Error_TransactionNotFound", func(t *testing.T) {
		nonExistentID := uuid.New()
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/transaction-history/%s", nonExistentID.String()), nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response handlers.DeleteTransactionResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response.Success)
		assert.Equal(t, "Transaction does not exist", response.Message)
	})

	// Test invalid transaction ID format
	t.Run("Error_InvalidTransactionID", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/transaction-history/invalid-id", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response handlers.DeleteTransactionResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response.Success)
		assert.Equal(t, "Invalid transaction ID format", response.Message)
	})
}

func TestDeleteTransactionsBatch(t *testing.T) {
	// Setup test handler and database
	transactionsHandler, db, err := setupTestTransactionsHandler(t)
	require.NoError(t, err)

	// Create test user
	user, err := createTestUser(db, "testuser2@example.com")
	require.NoError(t, err)

	// Create test transactions
	transactions := []models.Transaction{
		{
			UserID:          user.UserID,
			Symbol:          "AAPL",
			TradeType:       types.TradeType("Buy"),
			Quantity:        100,
			Price:           150.00,
			Amount:          15000.00,
			Currency:        "USD",
			Broker:          "Test Broker",
			Exchange:        "NASDAQ",
			TransactionDate: time.Now(),
			UserNotes:       "Test transaction 1",
		},
		{
			UserID:          user.UserID,
			Symbol:          "GOOGL",
			TradeType:       types.TradeType("Sell"),
			Quantity:        50,
			Price:           2500.00,
			Amount:          125000.00,
			Currency:        "USD",
			Broker:          "Test Broker",
			Exchange:        "NASDAQ",
			TransactionDate: time.Now(),
			UserNotes:       "Test transaction 2",
		},
	}

	for i := range transactions {
		err = db.Create(&transactions[i]).Error
		require.NoError(t, err)
	}

	// Setup gin context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/transaction-history", func(c *gin.Context) {
		c.Set("user_id", user.UserID)
		transactionsHandler.DeleteTransactions(c)
	})

	// Test successful batch deletion via request body
	t.Run("Success_BatchDeleteTransactions_RequestBody", func(t *testing.T) {
		requestBody := handlers.DeleteTransactionRequest{
			IDs: []string{
				transactions[0].TransactionID.String(),
				transactions[1].TransactionID.String(),
			},
		}

		body, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req, err := http.NewRequest("DELETE", "/transaction-history", bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.DeleteTransactionResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Equal(t, "2 transactions deleted successfully", response.Message)
		assert.Len(t, response.Data.DeletedIDs, 2)
	})

	// Create new transactions for query param test
	newTransactions := []models.Transaction{
		{
			UserID:          user.UserID,
			Symbol:          "TSLA",
			TradeType:       types.TradeType("Buy"),
			Quantity:        25,
			Price:           800.00,
			Amount:          20000.00,
			Currency:        "USD",
			Broker:          "Test Broker",
			Exchange:        "NASDAQ",
			TransactionDate: time.Now(),
			UserNotes:       "Test transaction 3",
		},
	}

	for i := range newTransactions {
		err = db.Create(&newTransactions[i]).Error
		require.NoError(t, err)
	}

	// Test successful batch deletion via query parameters
	t.Run("Success_BatchDeleteTransactions_QueryParams", func(t *testing.T) {
		url := fmt.Sprintf("/transaction-history?ids=%s", newTransactions[0].TransactionID.String())
		req, err := http.NewRequest("DELETE", url, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.DeleteTransactionResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Equal(t, "Transaction deleted successfully", response.Message)
		assert.Len(t, response.Data.DeletedIDs, 1)
	})

	// Test validation errors
	t.Run("Error_InvalidTransactionIDs", func(t *testing.T) {
		requestBody := handlers.DeleteTransactionRequest{
			IDs: []string{"invalid-id", "another-invalid-id"},
		}

		body, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req, err := http.NewRequest("DELETE", "/transaction-history", bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response handlers.DeleteTransactionResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response.Success)
		assert.Equal(t, "Validation failed", response.Message)
	})

	// Test empty request
	t.Run("Error_EmptyRequest", func(t *testing.T) {
		requestBody := handlers.DeleteTransactionRequest{
			IDs: []string{},
		}

		body, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req, err := http.NewRequest("DELETE", "/transaction-history", bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteTransactionOwnership(t *testing.T) {
	// Setup test handler and database
	transactionsHandler, db, err := setupTestTransactionsHandler(t)
	require.NoError(t, err)

	// Create two test users with unique usernames
	user1, err := createTestUserWithUsername(db, "user1_delete", "user1_delete@example.com")
	require.NoError(t, err)

	user2, err := createTestUserWithUsername(db, "user2_delete", "user2_delete@example.com")
	require.NoError(t, err)

	// Create a transaction for user1
	transaction := &models.Transaction{
		UserID:          user1.UserID,
		Symbol:          "AAPL",
		TradeType:       types.TradeType("Buy"),
		Quantity:        100,
		Price:           150.00,
		Amount:          15000.00,
		Currency:        "USD",
		Broker:          "Test Broker",
		Exchange:        "NASDAQ",
		TransactionDate: time.Now(),
		UserNotes:       "Test transaction for ownership test",
	}
	err = db.Create(transaction).Error
	require.NoError(t, err)

	// Setup gin context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/transaction-history/:id", func(c *gin.Context) {
		c.Set("user_id", user2.UserID) // Set user2 as the authenticated user
		transactionsHandler.DeleteTransaction(c)
	})

	// Test that user2 cannot delete user1's transaction
	t.Run("Error_Forbidden_NotOwner", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/transaction-history/%s", transaction.TransactionID.String()), nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code) // Should return not found for security

		var response handlers.DeleteTransactionResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response.Success)
		assert.Equal(t, "Transaction does not exist", response.Message)
	})
}
