package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"konkur-pulse/internal/handler"
	"konkur-pulse/internal/middleware"
	"konkur-pulse/internal/repository"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Database connection setup
	connStr := "postgres://pulse_user:pulse_password@localhost:5432/konkur_pulse?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cannot open database configuration: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	fmt.Println("✅ Successfully connected to Postgres!")

	// Initialize repository and database schema
	sessionRepo := repository.NewSessionRepository(db)
	if err := sessionRepo.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Initialize external API client
	mazeClient := repository.NewMazeClient("http://localhost:8081")

	// Setup HTTP router
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "ok", "message": "Konkur Pulse API is running", "db_connected": true}`)
	})

	// Register application endpoints
	mux.HandleFunc("GET /students/{id}/plan", handler.GetPlanHandler(mazeClient))
	mux.HandleFunc("POST /students/{id}/sessions", handler.SaveSessionHandler(sessionRepo))

	// Expose Prometheus metrics endpoint
	mux.Handle("GET /metrics", promhttp.Handler())

	// Configure structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := ":8080"
	logger.Info("Server is starting...", slog.String("port", port))

	// Wrap the router with observability middleware
	finalHandler := middleware.Observability(logger)(mux)

	if err := http.ListenAndServe(port, finalHandler); err != nil {
		logger.Error("Failed to start server", slog.String("error", err.Error()))
	}
}
