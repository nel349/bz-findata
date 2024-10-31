package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nel349/bz-findata/internal/analysis"
	"github.com/nel349/bz-findata/internal/analysis/database"
	"github.com/nel349/bz-findata/internal/analysis/supabase"
)

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

	// Example: Get largest orders in the last N hours
	hours := 24
	limit := 10
	largestOrders, err := analysisService.GetLargestReceivedOrdersInLastNHours(ctx, hours, limit)
	if err != nil {
		log.Fatalf("Failed to get largest orders: %v", err)
	}

	fmt.Printf("Largest `received` type orders in the last %d hours:\n", hours)
	for _, order := range largestOrders {
		localTime := time.Unix(order.Timestamp/1e9, 0).Local().UTC()
		location, err := time.LoadLocation("America/Denver")
		if err != nil {
			log.Fatalf("Failed to load location: %v", err)
		}
		fmt.Printf("Size: %f, Order ID: %s, Order Type: %s, Timestamp: %s, Product ID: %s, Price: %f\n",
			order.Size,
			order.OrderID,
			order.Type,
			localTime.In(location),
			order.ProductID,
			order.Price,
		)
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

	// get the largest received orders
	http.HandleFunc("/btc/largest-received-orders", func(w http.ResponseWriter, r *http.Request) {
		largestReceivedOrders, err := analysisService.GetLargestReceivedOrdersInLastNHours(ctx, 2, 10)
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
