package main

import (
	"context"
	"fmt"
	"log"

	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/provider"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create price service manager
	priceManager := provider.NewPriceServiceManager(cfg)

	// Test health check
	fmt.Println("=== Testing Health Check ===")
	health, err := priceManager.HealthCheck(context.Background())
	if err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Printf("Health: %+v\n", health)
	}

	// Test getting current prices
	fmt.Println("\n=== Testing Current Prices ===")
	symbols := []string{"AAPL", "GOOGL", "MSFT"}
	prices, err := priceManager.GetCurrentPrices(context.Background(), symbols)
	if err != nil {
		log.Printf("Failed to get current prices: %v", err)
	} else {
		for _, price := range prices {
			fmt.Printf("Symbol: %s, Price: $%.2f, Change: %.2f%%, Currency: %s\n",
				price.Symbol, price.CurrentPrice, price.ChangePercent, price.Currency)
		}
	}

	// Test getting single price
	fmt.Println("\n=== Testing Single Price ===")
	singlePrice, err := priceManager.GetCurrentPrice(context.Background(), "TSLA")
	if err != nil {
		log.Printf("Failed to get single price: %v", err)
	} else {
		fmt.Printf("TSLA Price: $%.2f, Change: %.2f%%\n",
			singlePrice.CurrentPrice, singlePrice.ChangePercent)
	}

	// Test historical prices
	fmt.Println("\n=== Testing Historical Prices ===")
	historicalPrices, err := priceManager.GetHistoricalPrices(
		context.Background(),
		[]string{"AAPL"},
		provider.ResolutionDaily,
		"2025-01-01",
		"2025-01-07",
	)
	if err != nil {
		log.Printf("Failed to get historical prices: %v", err)
	} else {
		for _, symbolData := range historicalPrices {
			fmt.Printf("Symbol: %s, Resolution: %s\n", symbolData.Symbol, symbolData.Resolution)
			for _, priceData := range symbolData.HistoricalPrices {
				fmt.Printf("  Date: %s, Price: $%.2f\n", priceData.Date, priceData.Price)
			}
		}
	}

	fmt.Println("\n=== Service Health Status ===")
	fmt.Printf("Is Service Healthy: %v\n", priceManager.IsServiceHealthy())
}
