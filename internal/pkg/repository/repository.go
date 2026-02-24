package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nekoteoj/lab-cms/internal/pkg/db"
)

// BaseRepository provides common functionality for all repositories.
type BaseRepository struct {
	dbManager *db.DBManager
	tableName string
}

// NewBaseRepository creates a new base repository.
func NewBaseRepository(dbManager *db.DBManager, tableName string) *BaseRepository {
	return &BaseRepository{
		dbManager: dbManager,
		tableName: tableName,
	}
}

// GetExecer returns the appropriate execer (DB or transaction) for the context.
func (r *BaseRepository) GetExecer(ctx context.Context) db.Execer {
	return r.dbManager.GetExecer(ctx)
}

// WithTransaction executes a function within a transaction.
func (r *BaseRepository) WithTransaction(ctx context.Context, fn db.TransactionFunc) error {
	return r.dbManager.WithTransaction(ctx, fn)
}

// CheckRowsAffected verifies that exactly one row was affected.
// Returns ErrNotFound if no rows were affected.
func CheckRowsAffected(result sql.Result, expected int64) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: failed to get rows affected: %v", ErrDatabase, err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	if affected != expected {
		return fmt.Errorf("%w: expected %d rows affected, got %d", ErrDatabase, expected, affected)
	}
	return nil
}
