package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nekoteoj/lab-cms/internal/pkg/config"
	"github.com/nekoteoj/lab-cms/internal/pkg/migrations"
	_ "modernc.org/sqlite"
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

	// Open database connection
	db, err := sql.Open("sqlite", cfg.DatabaseURL+"?_fk=1&_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Run migrations
	runner := migrations.NewRunner(db, "migrations")
	if err := runner.Run(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

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
