package database

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/transaction-tracker/backend/internal/models"
)

// Seeder handles database seeding operations
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// SeedDevelopmentData seeds the database with development data
func (s *Seeder) SeedDevelopmentData() error {
	log.Println("Seeding development data...")
	
	// Seed users
	if err := s.seedUsers(); err != nil {
		return err
	}
	
	// Seed transactions
	if err := s.seedTransactions(); err != nil {
		return err
	}
	
	log.Println("Development data seeded successfully")
	return nil
}

// seedUsers creates sample users
func (s *Seeder) seedUsers() error {
	// Check if users already exist
	var count int64
	if err := s.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		log.Println("Users already exist, skipping user seeding")
		return nil
	}
	
	log.Println("Seeding users...")
	
	// Hash password for demo users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	users := []models.User{
		{
			Username:     "demo_user",
			Email:        "demo@example.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "Demo",
			LastName:     "User",
			IsActive:     true,
		},
		{
			Username:     "john_doe",
			Email:        "john.doe@example.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "John",
			LastName:     "Doe",
			IsActive:     true,
		},
		{
			Username:     "jane_smith",
			Email:        "jane.smith@example.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "Jane",
			LastName:     "Smith",
			IsActive:     true,
		},
	}
	
	for _, user := range users {
		if err := s.db.Create(&user).Error; err != nil {
			return err
		}
		log.Printf("Created user: %s", user.Username)
	}
	
	return nil
}

// seedTransactions creates sample transactions
func (s *Seeder) seedTransactions() error {
	// Check if transactions already exist
	var count int64
	if err := s.db.Model(&models.Transaction{}).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		log.Println("Transactions already exist, skipping transaction seeding")
		return nil
	}
	
	log.Println("Seeding transactions...")
	
	// Get the demo user
	var demoUser models.User
	if err := s.db.Where("username = ?", "demo_user").First(&demoUser).Error; err != nil {
		return err
	}
	
	// Sample transaction data
	transactions := []models.Transaction{
		{
			UserID:          demoUser.ID,
			Type:            "buy",
			Symbol:          "AAPL",
			Quantity:        100,
			Price:           150.25,
			Amount:          15025.00,
			Fee:             1.00,
			Currency:        "USD",
			Broker:          "Firstrade",
			Account:         "Main",
			TransactionDate: time.Now().AddDate(0, 0, -30),
			Description:     "Buy Apple Inc. shares",
			Status:          "completed",
			ExtractedFrom:   "manual",
			Metadata:        "{}",
		},
		{
			UserID:          demoUser.ID,
			Type:            "buy",
			Symbol:          "GOOGL",
			Quantity:        50,
			Price:           2800.00,
			Amount:          140000.00,
			Fee:             1.00,
			Currency:        "USD",
			Broker:          "Firstrade",
			Account:         "Main",
			TransactionDate: time.Now().AddDate(0, 0, -25),
			Description:     "Buy Alphabet Inc. shares",
			Status:          "completed",
			ExtractedFrom:   "manual",
			Metadata:        "{}",
		},
		{
			UserID:          demoUser.ID,
			Type:            "sell",
			Symbol:          "AAPL",
			Quantity:        50,
			Price:           155.75,
			Amount:          7787.50,
			Fee:             1.00,
			Currency:        "USD",
			Broker:          "Firstrade",
			Account:         "Main",
			TransactionDate: time.Now().AddDate(0, 0, -20),
			Description:     "Sell Apple Inc. shares",
			Status:          "completed",
			ExtractedFrom:   "manual",
			Metadata:        "{}",
		},
		{
			UserID:          demoUser.ID,
			Type:            "buy",
			Symbol:          "MSFT",
			Quantity:        75,
			Price:           415.30,
			Amount:          31147.50,
			Fee:             1.00,
			Currency:        "USD",
			Broker:          "Firstrade",
			Account:         "Main",
			TransactionDate: time.Now().AddDate(0, 0, -15),
			Description:     "Buy Microsoft Corp. shares",
			Status:          "completed",
			ExtractedFrom:   "manual",
			Metadata:        "{}",
		},
		{
			UserID:          demoUser.ID,
			Type:            "dividend",
			Symbol:          "AAPL",
			Quantity:        50,
			Price:           0.24,
			Amount:          12.00,
			Fee:             0.00,
			Currency:        "USD",
			Broker:          "Firstrade",
			Account:         "Main",
			TransactionDate: time.Now().AddDate(0, 0, -10),
			Description:     "Apple Inc. dividend payment",
			Status:          "completed",
			ExtractedFrom:   "manual",
			Metadata:        "{}",
		},
		{
			UserID:          demoUser.ID,
			Type:            "buy",
			Symbol:          "TSLA",
			Quantity:        25,
			Price:           245.80,
			Amount:          6145.00,
			Fee:             1.00,
			Currency:        "USD",
			Broker:          "Firstrade",
			Account:         "Main",
			TransactionDate: time.Now().AddDate(0, 0, -5),
			Description:     "Buy Tesla Inc. shares",
			Status:          "completed",
			ExtractedFrom:   "manual",
			Metadata:        "{}",
		},
		{
			UserID:          demoUser.ID,
			Type:            "buy",
			Symbol:          "NVDA",
			Quantity:        10,
			Price:           875.50,
			Amount:          8755.00,
			Fee:             1.00,
			Currency:        "USD",
			Broker:          "Firstrade",
			Account:         "Main",
			TransactionDate: time.Now().AddDate(0, 0, -3),
			Description:     "Buy NVIDIA Corp. shares",
			Status:          "completed",
			ExtractedFrom:   "manual",
			Metadata:        "{}",
		},
	}
	
	// Get the john_doe user for additional transactions
	var johnUser models.User
	if err := s.db.Where("username = ?", "john_doe").First(&johnUser).Error; err != nil {
		return err
	}
	
	// Add some transactions for john_doe
	johnTransactions := []models.Transaction{
		{
			UserID:          johnUser.ID,
			Type:            "buy",
			Symbol:          "SPY",
			Quantity:        200,
			Price:           420.50,
			Amount:          84100.00,
			Fee:             0.00,
			Currency:        "USD",
			Broker:          "Vanguard",
			Account:         "IRA",
			TransactionDate: time.Now().AddDate(0, 0, -45),
			Description:     "Buy SPDR S&P 500 ETF",
			Status:          "completed",
			ExtractedFrom:   "manual",
			Metadata:        "{}",
		},
		{
			UserID:          johnUser.ID,
			Type:            "buy",
			Symbol:          "VTI",
			Quantity:        100,
			Price:           225.75,
			Amount:          22575.00,
			Fee:             0.00,
			Currency:        "USD",
			Broker:          "Vanguard",
			Account:         "Taxable",
			TransactionDate: time.Now().AddDate(0, 0, -40),
			Description:     "Buy Vanguard Total Stock Market ETF",
			Status:          "completed",
			ExtractedFrom:   "manual",
			Metadata:        "{}",
		},
	}
	
	transactions = append(transactions, johnTransactions...)
	
	// Create transactions in batches
	if err := s.db.CreateInBatches(transactions, 10).Error; err != nil {
		return err
	}
	
	log.Printf("Created %d sample transactions", len(transactions))
	return nil
}

// SeedTestData seeds the database with test data
func (s *Seeder) SeedTestData() error {
	log.Println("Seeding test data...")
	
	// Create a test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	testUser := models.User{
		Username:     "test_user",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FirstName:    "Test",
		LastName:     "User",
		IsActive:     true,
	}
	
	if err := s.db.Create(&testUser).Error; err != nil {
		return err
	}
	
	// Create test transactions
	testTransaction := models.Transaction{
		UserID:          testUser.ID,
		Type:            "buy",
		Symbol:          "TEST",
		Quantity:        1,
		Price:           1.00,
		Amount:          1.00,
		Fee:             0.00,
		Currency:        "USD",
		Broker:          "Test Broker",
		Account:         "Test",
		TransactionDate: time.Now(),
		Description:     "Test transaction",
		Status:          "completed",
		ExtractedFrom:   "test",
		Metadata:        "{}",
	}
	
	if err := s.db.Create(&testTransaction).Error; err != nil {
		return err
	}
	
	log.Println("Test data seeded successfully")
	return nil
}

// ClearTestData removes all test data
func (s *Seeder) ClearTestData() error {
	log.Println("Clearing test data...")
	
	// Delete test transactions
	if err := s.db.Unscoped().Where("extracted_from = ?", "test").Delete(&models.Transaction{}).Error; err != nil {
		return err
	}
	
	// Delete test user
	if err := s.db.Unscoped().Where("username = ?", "test_user").Delete(&models.User{}).Error; err != nil {
		return err
	}
	
	log.Println("Test data cleared successfully")
	return nil
}
