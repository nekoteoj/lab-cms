// Package migrations provides database migration functionality for the Lab CMS.
// It manages schema versioning and applies migrations in sequential order.
package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// Migration represents a single database migration.
type Migration struct {
	Version int
	Name    string
	SQL     string
}

// Runner manages database migrations.
type Runner struct {
	db            *sql.DB
	migrationsDir string
}

// NewRunner creates a new migration runner.
// It takes a database connection and the path to the migrations directory.
func NewRunner(db *sql.DB, migrationsDir string) *Runner {
	return &Runner{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Run applies all pending migrations.
// It creates the schema_migrations table if it doesn't exist,
// reads migration files from the migrations directory,
// and applies any migrations that haven't been applied yet.
func (r *Runner) Run() error {
	// Enable foreign keys
	if _, err := r.db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := r.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations, err := r.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	if len(migrations) == 0 {
		return nil
	}

	applied, err := r.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, migration := range migrations {
		if _, ok := applied[migration.Version]; ok {
			continue
		}

		if err := r.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %d (%s): %w",
				migration.Version, migration.Name, err)
		}
	}

	return nil
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist.
func (r *Runner) createMigrationsTable() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// loadMigrations reads migration files from the migrations directory.
// It returns migrations sorted by version number.
func (r *Runner) loadMigrations() ([]Migration, error) {
	files, err := os.ReadDir(r.migrationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		migration, err := r.parseMigrationFile(file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to parse migration file %s: %w", file.Name(), err)
		}

		migrations = append(migrations, migration)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// parseMigrationFile extracts version and name from a migration filename.
// Expected format: 001_migration_name.sql
func (r *Runner) parseMigrationFile(filename string) (Migration, error) {
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) != 2 {
		return Migration{}, fmt.Errorf("invalid migration filename format: %s", filename)
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return Migration{}, fmt.Errorf("invalid version number in filename %s: %w", filename, err)
	}

	name := strings.TrimSuffix(parts[1], ".sql")

	content, err := os.ReadFile(filepath.Join(r.migrationsDir, filename))
	if err != nil {
		return Migration{}, fmt.Errorf("failed to read migration file %s: %w", filename, err)
	}

	return Migration{
		Version: version,
		Name:    name,
		SQL:     string(content),
	}, nil
}

// getAppliedMigrations returns a map of already applied migration versions.
func (r *Runner) getAppliedMigrations() (map[int]bool, error) {
	rows, err := r.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// applyMigration executes a single migration within a transaction.
func (r *Runner) applyMigration(m Migration) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(m.SQL); err != nil {
		return fmt.Errorf("migration SQL failed: %w", err)
	}

	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (version) VALUES (?)",
		m.Version,
	); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetPendingMigrations returns a list of migration versions that haven't been applied yet.
func (r *Runner) GetPendingMigrations() ([]int, error) {
	migrations, err := r.loadMigrations()
	if err != nil {
		return nil, err
	}

	applied, err := r.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	var pending []int
	for _, m := range migrations {
		if _, ok := applied[m.Version]; !ok {
			pending = append(pending, m.Version)
		}
	}

	return pending, nil
}

// GetAppliedMigrations returns a list of already applied migration versions.
func (r *Runner) GetAppliedMigrations() ([]int, error) {
	rows, err := r.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, rows.Err()
}
