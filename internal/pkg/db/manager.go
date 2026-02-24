// Package db provides database connection management and transaction support.
// The DBManager wraps sql.DB (which is already a connection pool) and provides
// helper methods for common operations.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const txContextKey contextKey = "db_transaction"

// DBManager wraps sql.DB to provide a unified interface for database operations.
// sql.DB is already a connection pool safe for concurrent use across goroutines.
type DBManager struct {
	db *sql.DB
}

// NewManager creates a new DBManager with the given database URL.
// The database is opened with WAL mode and foreign key constraints enabled.
func NewManager(databaseURL string) (*DBManager, error) {
	// Open database with WAL mode and foreign key constraints
	db, err := sql.Open("sqlite", databaseURL+"?_fk=1&_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DBManager{db: db}, nil
}

// ConfigurePool sets the connection pool limits.
// Pass 0 for maxOpenConns or maxIdleConns to use Go defaults.
func (m *DBManager) ConfigurePool(maxOpenConns, maxIdleConns int) {
	if maxOpenConns > 0 {
		m.db.SetMaxOpenConns(maxOpenConns)
		log.Printf("Database pool: max open connections set to %d", maxOpenConns)
	}
	if maxIdleConns > 0 {
		m.db.SetMaxIdleConns(maxIdleConns)
		log.Printf("Database pool: max idle connections set to %d", maxIdleConns)
	}
}

// GetDB returns the underlying sql.DB instance.
// Use this for direct database access when needed.
func (m *DBManager) GetDB() *sql.DB {
	return m.db
}

// Ping checks if the database connection is alive.
func (m *DBManager) Ping(ctx context.Context) error {
	return m.db.PingContext(ctx)
}

// Close closes the database connection pool.
// Should be called during graceful shutdown.
func (m *DBManager) Close() error {
	return m.db.Close()
}

// TransactionFunc is a function that executes within a transaction.
// The function receives a context with the transaction stored in it.
// If the function returns an error, the transaction is rolled back.
// If the function returns nil, the transaction is committed.
type TransactionFunc func(ctx context.Context) error

// WithTransaction executes the given function within a database transaction.
// The transaction is committed if the function returns nil, otherwise it's rolled back.
// The transaction is stored in the context and can be retrieved using GetTx(ctx).
func (m *DBManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Store transaction in context
	txCtx := context.WithValue(ctx, txContextKey, tx)

	// Execute the function
	if err := fn(txCtx); err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %w; rollback also failed: %v", err, rbErr)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTx retrieves the transaction from the context.
// Returns nil if no transaction is in the context.
func GetTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txContextKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}

// Execer is an interface that can execute SQL statements.
// Both *sql.DB and *sql.Tx implement this interface.
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// GetExecer returns an Execer for the given context.
// If a transaction is present in the context, it returns the transaction.
// Otherwise, it returns the database connection.
func (m *DBManager) GetExecer(ctx context.Context) Execer {
	if tx := GetTx(ctx); tx != nil {
		return tx
	}
	return m.db
}
