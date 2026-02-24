package test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/nekoteoj/lab-cms/internal/pkg/migrations"
	"github.com/nekoteoj/lab-cms/test/helpers"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestRunner_CreatesMigrationsTable(t *testing.T) {
	db := helpers.NewTestDB(t)

	exists := helpers.TableExists(t, db, "schema_migrations")
	require.True(t, exists, "schema_migrations table should exist")
}

func TestRunner_RunsMigrationsInOrder(t *testing.T) {
	db := helpers.NewTestDB(t)

	// Verify tables from both migrations exist
	tables := []string{
		"users",
		"lab_members",
		"publications",
		"projects",
		"news",
		"homepage_sections",
		"lab_settings",
		"project_members",
		"publication_authors",
		"project_publications",
	}

	for _, table := range tables {
		exists := helpers.TableExists(t, db, table)
		require.True(t, exists, "table %s should exist", table)
	}
}

func TestRunner_TracksAppliedMigrations(t *testing.T) {
	db := helpers.NewTestDB(t)

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	require.NoError(t, err)
	require.GreaterOrEqual(t, count, 2, "should have at least 2 migrations recorded")
}

func TestRunner_IdempotentExecution(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_fk=1")
	require.NoError(t, err)
	defer db.Close()

	runner := migrations.NewRunner(db, "../migrations")

	// First run
	err = runner.Run()
	require.NoError(t, err)

	var firstCount int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&firstCount)
	require.NoError(t, err)

	// Second run (should be idempotent)
	err = runner.Run()
	require.NoError(t, err)

	var secondCount int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&secondCount)
	require.NoError(t, err)

	require.Equal(t, firstCount, secondCount, "running migrations twice should not add duplicate entries")
}

func TestRunner_GetAppliedMigrations(t *testing.T) {
	db := helpers.NewTestDB(t)

	runner := migrations.NewRunner(db, "../migrations")
	applied, err := runner.GetAppliedMigrations()
	require.NoError(t, err)
	require.NotEmpty(t, applied)

	// Verify migrations are in ascending order
	for i := 1; i < len(applied); i++ {
		require.Less(t, applied[i-1], applied[i], "migrations should be in ascending order")
	}
}

func TestRunner_GetPendingMigrations(t *testing.T) {
	db := helpers.NewTestDB(t)

	runner := migrations.NewRunner(db, "../migrations")
	pending, err := runner.GetPendingMigrations()
	require.NoError(t, err)
	require.Empty(t, pending, "should have no pending migrations after running all")
}

func TestRunner_NoMigrationsDirectory(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_fk=1")
	require.NoError(t, err)
	defer db.Close()

	runner := migrations.NewRunner(db, "./nonexistent_migrations")
	err = runner.Run()
	require.NoError(t, err, "should not error when migrations directory doesn't exist")
}

func TestRunner_InvalidMigrationFile(t *testing.T) {
	// Create temporary directory with invalid migration file
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "invalid.sql"), []byte("SELECT 1"), 0644)
	require.NoError(t, err)

	db, err := sql.Open("sqlite", ":memory:?_fk=1")
	require.NoError(t, err)
	defer db.Close()

	runner := migrations.NewRunner(db, tmpDir)
	err = runner.Run()
	require.Error(t, err, "should error on invalid migration filename")
}

func TestRunner_MigrationSQLFailure(t *testing.T) {
	// Create temporary directory with a migration that will fail
	tmpDir := t.TempDir()
	err := os.WriteFile(
		filepath.Join(tmpDir, "001_broken.sql"),
		[]byte("THIS IS NOT VALID SQL"),
		0644,
	)
	require.NoError(t, err)

	db, err := sql.Open("sqlite", ":memory:?_fk=1")
	require.NoError(t, err)
	defer db.Close()

	runner := migrations.NewRunner(db, tmpDir)
	err = runner.Run()
	require.Error(t, err, "should error when migration SQL fails")
	require.Contains(t, err.Error(), "broken", "error should reference the failing migration")
}

func TestRunner_RollbackOnFailure(t *testing.T) {
	// Create temporary directory with a good first migration and a bad second
	tmpDir := t.TempDir()

	// Good first migration
	err := os.WriteFile(
		filepath.Join(tmpDir, "001_good.sql"),
		[]byte("CREATE TABLE test_table (id INTEGER PRIMARY KEY)"),
		0644,
	)
	require.NoError(t, err)

	// Bad second migration
	err = os.WriteFile(
		filepath.Join(tmpDir, "001_broken.sql"),
		[]byte("THIS IS NOT VALID SQL SYNTAX AT ALL"),
		0644,
	)
	require.NoError(t, err)

	db, err := sql.Open("sqlite", ":memory:?_fk=1")
	require.NoError(t, err)
	defer db.Close()

	runner := migrations.NewRunner(db, tmpDir)
	err = runner.Run()
	require.Error(t, err)

	// Verify first migration was rolled back (table shouldn't exist)
	var exists bool
	err = db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='test_table')",
	).Scan(&exists)
	require.NoError(t, err)
	require.False(t, exists, "first migration should have been rolled back")

	// Verify second migration was not recorded
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = 2").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count, "second migration should not be recorded")
}
