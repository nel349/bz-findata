package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nel349/coinbase-analysis/auth"
	"github.com/nel349/coinbase-analysis/internal/analysis"
	"github.com/nel349/coinbase-analysis/internal/database"
)

const (
	requestMethod = "GET"
	requestHost   = "api.coinbase.com"
	requestPath   = "/api/v3/brokerage/accounts"
)

func main() {
	ctx := context.Background()

	// Initialize database connection
	db, err := database.NewConnection(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_BASE"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize analysis service
	analysisService := analysis.NewService(db)

	fmt.Println("Starting analysis app...")

	// Example: Get largest orders in the last 24 hours
	largestOrders, err := analysisService.GetLargestOrdersInLastNHours(ctx, 24, 10)
	if err != nil {
		log.Fatalf("Failed to get largest orders: %v", err)
	}

	fmt.Println("Largest orders in the last 24 hours:")
	for _, order := range largestOrders {
		fmt.Printf("Order ID: %s, Size: %f, Price: %f\n", order.OrderID, order.Size, order.Price)
	}

	uri := fmt.Sprintf("%s %s%s", requestMethod, requestHost, requestPath)

	jwt, err := auth.BuildJWT(uri)

	fmt.Println("adslfdksfjdsfsdf", jwt)

	if err != nil {
		fmt.Printf("error building jwt: %v", err)
	}
}