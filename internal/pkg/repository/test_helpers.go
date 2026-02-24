package repository

import (
	"context"
	"testing"

	"github.com/nekoteoj/lab-cms/internal/pkg/db"
	"github.com/nekoteoj/lab-cms/internal/pkg/migrations"
	"github.com/stretchr/testify/require"
)

// ctx is the shared background context for all tests
var ctx = context.Background()

// setupTestDB creates a test database with migrations for repository tests
func setupTestDB(t *testing.T) *db.DBManager {
	dbManager, err := db.NewManager(":memory:")
	require.NoError(t, err)

	t.Cleanup(func() {
		dbManager.Close()
	})

	runner := migrations.NewRunner(dbManager.GetDB(), "../../../migrations")
	err = runner.Run()
	require.NoError(t, err)

	return dbManager
}
