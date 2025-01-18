package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nel349/bz-findata/config"
	"github.com/nel349/bz-findata/internal/analysis/database"
	"github.com/nel349/bz-findata/internal/analysis/dex"
	"github.com/nel349/bz-findata/internal/analysis/handlers"
	"github.com/nel349/bz-findata/internal/analysis/infrastructure/scheduler"
	"github.com/nel349/bz-findata/internal/analysis/orders"
	"github.com/nel349/bz-findata/internal/analysis/supabase"
	"github.com/nel349/bz-findata/internal/analysis/task"
	"github.com/robfig/cron/v3"
)

type ScheduledTask struct {
	ID       cron.EntryID `json:"id"`
	Schedule string       `json:"schedule"`
	Hours    int          `json:"hours"`
	Limit    int          `json:"limit"`
}

func main() {
	// var dbPassword string

	// Check if we're running locally (you can set this env var in docker-compose)
	// if os.Getenv("IS_LOCAL") == "true" {
	//     dbPassword = os.Getenv("DB_PASSWORD")
	// } else {
	//     // Get from AWS Secrets Manager
	//     dbSecret, err := awslocal.GetDefaultDBSecret()
	//     if err != nil {
	//         log.Fatalf("Failed to retrieve DB secret: %v", err)
	//     }
	//     dbPassword = dbSecret.DB_PASSWORD

	// 	// log.Println("DB Password:", dbPassword) // TODO: remove
	// 	// // kill process
	// 	// os.Exit(0)
	// }

	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	cfg, err := config.NewAnalysisConfig(ctx)
	if err != nil {
		log.Fatalf("failed config init: %v", err)
	}

	// Setup dependencies
	if err := run(cfg); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}

func run(cfg *config.AnalysisConfig) error {
	// Initialize dependencies
	db, err := database.NewConnection(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Base,
	)
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialize services
	supabaseRepo := supabase.NewSupabaseRepository()
	analysisService := analysis.NewService(db, supabaseRepo.Client)
	dexService := dex.NewService(db, supabaseRepo.Client)
	
	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(analysisService)
	dexHandler := handlers.NewDexHandler(dexService)


	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Initialize task manager
	taskService := task.NewService(analysisService)
	taskManager := scheduler.NewTaskManager(taskService)

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/btc", func(r chi.Router) {

			// Largest orders and store in supabase
			r.Get("/largest-received-orders", orderHandler.GetLargestReceivedOrders)
			r.Get("/largest-open-orders", orderHandler.GetLargestOpenOrders)
			r.Get("/largest-match-orders", orderHandler.GetLargestMatchOrders)
			r.Post("/store-received-orders", orderHandler.StoreReceivedOrdersInSupabase)
			r.Post("/store-match-orders", orderHandler.StoreMatchOrdersInSupabase)

			// Scheduler endpoints
			r.Route("/scheduler", func(r chi.Router) {
				r.Post("/start", taskManager.StartTask)
				r.Delete("/stop/{taskID}", taskManager.StopTask)
				r.Get("/tasks", taskManager.ListTasks)
			})
		})

		// Routes for dex
		r.Route("/dex", func(r chi.Router) {
			r.Get("/largest-swaps", dexHandler.GetLargestSwaps)
		})
	})

	log.Printf("Server started on port %s", "8090")
	return http.ListenAndServe(":8090", r)
}
