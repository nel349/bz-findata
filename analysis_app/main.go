package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// "github.com/nel349/coinbase-analysis/auth"
	// "github.com/nel349/bz-findata/pkg/exchange/coinbase"
	"github.com/nel349/coinbase-analysis/internal/analysis"
	"github.com/nel349/coinbase-analysis/internal/database"
	"github.com/nel349/coinbase-analysis/supabase"
)

// const (
// 	requestMethod = "GET"
// 	requestHost   = "api.coinbase.com"
// 	requestPath   = "/api/v3/brokerage/key_permissions"
// )

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

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

	// Initialize supabase client
	supabaseRepo := supabase.NewSupabaseRepository()
	// Initialize analysis service
	analysisService := analysis.NewService(db, supabaseRepo.Client)

	fmt.Println("Starting analysis app...")

	// Example: Get largest orders in the last 24 hours
	largestOrders, err := analysisService.GetLargestOrdersInLastNHours(ctx, 24, 10)
	if err != nil {
		log.Fatalf("Failed to get largest orders: %v", err)
	}

	fmt.Println("Largest orders in the last 24 hours:")
	for _, order := range largestOrders {
		fmt.Printf("Order ID: %s, Product ID: %s, Price: %f\n", order.Type, order.ProductID, order.Price)
	}

	// uri := fmt.Sprintf("https://%s%s", requestHost, requestPath)

	// jwt, err := auth.BuildJWT(uri)

	// fmt.Println("JWT:", jwt)
	// fmt.Println("URI:", uri)

	// lets print the command with the jwt
	// fmt.Printf("export JWT=%s\n", jwt)

	// print curl command
	// fmt.Printf("curl -X %s https://%s%s -H \"Authorization: Bearer %s\"\n", requestMethod, requestHost, requestPath, jwt)

	// auth := coinbase.NewAuth()
	// 	auth.GenerateSignature()

	if err != nil {
		fmt.Printf("error building jwt: %v", err)
	}

	// Message that the server is running
	fmt.Println("Server is running on port 8090")

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	// get the largest orders
	http.HandleFunc("/btc/largest-orders", func(w http.ResponseWriter, r *http.Request) {
		largestOrders, err := analysisService.GetLargestOrdersInLastNHours(ctx, 24, 10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(largestOrders)
	})

	// get the largest received orders
	http.HandleFunc("/btc/largest-received-orders", func(w http.ResponseWriter, r *http.Request) {
		largestReceivedOrders, err := analysisService.GetLargestReceivedOrdersInLastNHours(ctx, 24, 10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(largestReceivedOrders)
	})

	// get the largest open orders
	http.HandleFunc("/btc/largest-open-orders", func(w http.ResponseWriter, r *http.Request) {
		largestOpenOrders, err := analysisService.GetLargestOpenOrdersInLastNHours(ctx, 24, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(largestOpenOrders)
	})

	http.ListenAndServe(":8090", nil)
}
