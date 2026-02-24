package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	t.Run("create in-memory database", func(t *testing.T) {
		dbManager, err := NewManager(":memory:")
		require.NoError(t, err)
		defer dbManager.Close()

		assert.NotNil(t, dbManager)
		assert.NotNil(t, dbManager.GetDB())
	})

	t.Run("database is accessible", func(t *testing.T) {
		dbManager, err := NewManager(":memory:")
		require.NoError(t, err)
		defer dbManager.Close()

		// Ping should succeed
		err = dbManager.Ping(context.Background())
		require.NoError(t, err)
	})
}

func TestDBManager_ConfigurePool(t *testing.T) {
	dbManager, err := NewManager(":memory:")
	require.NoError(t, err)
	defer dbManager.Close()

	t.Run("configure max open connections", func(t *testing.T) {
		dbManager.ConfigurePool(10, 0)
		// Configuration applied without error
		assert.Equal(t, 10, dbManager.GetDB().Stats().MaxOpenConnections)
	})

	t.Run("configure max idle connections", func(t *testing.T) {
		dbManager.ConfigurePool(10, 5)
		// Configuration applied without error
		// Note: idle connections may be 0 if not actively used
	})

	t.Run("zero values use defaults", func(t *testing.T) {
		dbManager.ConfigurePool(0, 0)
		// Should not panic and use Go defaults
	})
}

func TestDBManager_WithTransaction(t *testing.T) {
	dbManager, err := NewManager(":memory:")
	require.NoError(t, err)
	defer dbManager.Close()

	// Create a test table
	db := dbManager.GetDB()
	_, err = db.Exec(`
		CREATE TABLE test_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	t.Run("successful transaction commits", func(t *testing.T) {
		ctx := context.Background()
		err := dbManager.WithTransaction(ctx, func(txCtx context.Context) error {
			tx := GetTx(txCtx)
			require.NotNil(t, tx)

			_, err := tx.ExecContext(txCtx, "INSERT INTO test_items (name) VALUES (?)", "test")
			return err
		})
		require.NoError(t, err)

		// Verify data was committed
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM test_items WHERE name = ?", "test").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("transaction with error rolls back", func(t *testing.T) {
		ctx := context.Background()
		testErr := assert.AnError
		err := dbManager.WithTransaction(ctx, func(txCtx context.Context) error {
			tx := GetTx(txCtx)
			require.NotNil(t, tx)

			_, err := tx.ExecContext(txCtx, "INSERT INTO test_items (name) VALUES (?)", "rollback_test")
			require.NoError(t, err)

			// Return error to trigger rollback
			return testErr
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, testErr)

		// Verify data was NOT committed (rolled back)
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM test_items WHERE name = ?", "rollback_test").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("nested transaction uses same outer transaction", func(t *testing.T) {
		ctx := context.Background()
		err := dbManager.WithTransaction(ctx, func(txCtx context.Context) error {
			outerTx := GetTx(txCtx)
			require.NotNil(t, outerTx)

			// Insert in outer transaction
			_, err := outerTx.ExecContext(txCtx, "INSERT INTO test_items (name) VALUES (?)", "nested_test")
			require.NoError(t, err)

			// Try to start another transaction - should use savepoint in SQLite
			err = dbManager.WithTransaction(txCtx, func(innerCtx context.Context) error {
				innerTx := GetTx(innerCtx)
				// In SQLite, nested transactions create savepoints but GetTx should return a transaction
				assert.NotNil(t, innerTx)
				// The inner context should have transaction from context
				return nil
			})
			require.NoError(t, err)

			return nil
		})
		require.NoError(t, err)

		// Verify the data was committed
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM test_items WHERE name = ?", "nested_test").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestGetTx(t *testing.T) {
	t.Run("returns nil when no transaction", func(t *testing.T) {
		ctx := context.Background()
		tx := GetTx(ctx)
		assert.Nil(t, tx)
	})

	t.Run("returns transaction when in context", func(t *testing.T) {
		dbManager, err := NewManager(":memory:")
		require.NoError(t, err)
		defer dbManager.Close()

		ctx := context.Background()
		err = dbManager.WithTransaction(ctx, func(txCtx context.Context) error {
			tx := GetTx(txCtx)
			assert.NotNil(t, tx)
			return nil
		})
		require.NoError(t, err)
	})
}

func TestDBManager_GetExecer(t *testing.T) {
	dbManager, err := NewManager(":memory:")
	require.NoError(t, err)
	defer dbManager.Close()

	t.Run("returns DB when no transaction", func(t *testing.T) {
		ctx := context.Background()
		execer := dbManager.GetExecer(ctx)
		assert.NotNil(t, execer)
		// Should be the underlying *sql.DB
	})

	t.Run("returns Tx when in transaction", func(t *testing.T) {
		ctx := context.Background()
		err := dbManager.WithTransaction(ctx, func(txCtx context.Context) error {
			execer := dbManager.GetExecer(txCtx)
			assert.NotNil(t, execer)
			// Should be the *sql.Tx from context
			return nil
		})
		require.NoError(t, err)
	})
}
