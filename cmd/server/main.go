package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nekoteoj/lab-cms/internal/pkg/config"
	"github.com/nekoteoj/lab-cms/internal/pkg/db"
	"github.com/nekoteoj/lab-cms/internal/pkg/migrations"
	"github.com/nekoteoj/lab-cms/internal/pkg/repository"
)

func main() {
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

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
	log.Println("Database migrations completed successfully")

	// Initialize repository factory
	repoFactory := repository.NewFactory(dbManager)
	// Note: repoFactory will be used by HTTP handlers in future tasks
	_ = repoFactory // Prevent unused variable warning

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Lab CMS")
	})

	log.Printf("Server starting on port %s [%s]", cfg.Port, cfg.Env)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// ensureDataDir creates the parent directory for the database file if it doesn't exist.
func ensureDataDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir == "." || dir == "/" {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}
