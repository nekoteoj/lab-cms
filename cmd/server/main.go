package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/nekoteoj/lab-cms/internal/app/server"
	"github.com/nekoteoj/lab-cms/internal/pkg/config"
	"github.com/nekoteoj/lab-cms/internal/pkg/db"
	"github.com/nekoteoj/lab-cms/internal/pkg/logger"
	"github.com/nekoteoj/lab-cms/internal/pkg/migrations"
	"github.com/nekoteoj/lab-cms/internal/pkg/repository"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Init("error", cfg.IsProduction())
		logger.L().Fatal("Configuration error: " + err.Error())
	}

	// Initialize logger with configuration
	logger.Init(cfg.LogLevel, cfg.IsProduction())
	log := logger.L()

	log.Info("Starting Lab CMS")
	log.WithField("port", cfg.Port).
		WithField("env", cfg.Env).
		Info("Configuration loaded")

	// Ensure data directory exists
	if err := ensureDataDir(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize database manager with connection pool
	dbManager, err := db.NewManager(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbManager.Close()

	// Configure connection pool (optional, uses Go defaults if 0)
	dbManager.ConfigurePool(cfg.DBMaxOpenConns, cfg.DBMaxIdleConns)

	// Run migrations
	runner := migrations.NewRunner(dbManager.GetDB(), "migrations")
	if err := runner.Run(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Info("Database migrations completed successfully")

	// Initialize repository factory
	repoFactory := repository.NewFactory(dbManager)
	// Note: repoFactory will be used by HTTP handlers in future tasks
	_ = repoFactory // Prevent unused variable warning

	// Set up HTTP handlers with middleware chain
	handler := setupHandler(cfg)

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithField("address", srv.Addr).Info("Server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutdown signal received, gracefully shutting down...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}

// setupHandler creates the HTTP handler with middleware chain
func setupHandler(cfg *config.Config) http.Handler {
	// Create base mux
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	// Home route (placeholder)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			server.RespondNotFound(w, r, "page")
			return
		}
		fmt.Fprintf(w, "Welcome to Lab CMS")
	})

	// Apply middleware chain
	middlewares := []server.Middleware{
		server.RequestIDMiddleware(),
		server.RecoveryMiddleware(),
		server.SecurityHeadersMiddleware(),
		server.LoggingMiddleware(),
	}

	return server.Chain(middlewares...)(mux)
}

// ensureDataDir creates the parent directory for the database file if it doesn't exist.
func ensureDataDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir == "." || dir == "/" {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}
