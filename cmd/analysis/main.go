package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nel349/bz-findata/internal/analysis"
	"github.com/nel349/bz-findata/internal/analysis/database"
	"github.com/nel349/bz-findata/internal/analysis/handlers"
	"github.com/nel349/bz-findata/internal/analysis/supabase"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
}

func main() {
	// Initialize configuration
	cfg := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_BASE"),
		Port:       "8090",
	}

	// Setup dependencies
	if err := run(cfg); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}

func run(cfg Config) error {
	// Initialize dependencies
	db, err := database.NewConnection(
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialize services
	supabaseRepo := supabase.NewSupabaseRepository()
	analysisService := analysis.NewService(db, supabaseRepo.Client)

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(analysisService)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/btc", func(r chi.Router) {
			r.Get("/largest-received-orders", orderHandler.GetLargestReceivedOrders)
			r.Get("/largest-open-orders", orderHandler.GetLargestOpenOrders)
			r.Get("/largest-match-orders", orderHandler.GetLargestMatchOrders)
			r.Post("/store-received-orders", orderHandler.StoreReceivedOrdersInSupabase)
			r.Post("/store-match-orders", orderHandler.StoreMatchOrdersInSupabase)
		})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	return http.ListenAndServe(":"+cfg.Port, r)
}