// Package helpers provides utilities for database testing.
package helpers

import (
	"database/sql"
	"testing"

	"github.com/nekoteoj/lab-cms/internal/pkg/migrations"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// NewTestDB creates a new in-memory SQLite database for testing.
// It enables foreign keys and runs all migrations automatically.
// Returns the database connection which will be cleaned up when the test completes.
func NewTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err, "failed to open in-memory database")

	t.Cleanup(func() {
		db.Close()
	})

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err, "failed to enable foreign keys")

	runner := migrations.NewRunner(db, "../migrations")
	err = runner.Run()
	require.NoError(t, err, "failed to run migrations")

	return db
}

// NewTestDBWithMigrations creates a test DB but allows specifying a custom migrations directory.
// Useful for testing specific migration scenarios.
func NewTestDBWithMigrations(t *testing.T, migrationsDir string) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err, "failed to open in-memory database")

	t.Cleanup(func() {
		db.Close()
	})

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err, "failed to enable foreign keys")

	runner := migrations.NewRunner(db, migrationsDir)
	err = runner.Run()
	require.NoError(t, err, "failed to run migrations")

	return db
}

// TableExists checks if a table exists in the database.
func TableExists(t *testing.T, db *sql.DB, tableName string) bool {
	var exists bool
	err := db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM sqlite_master 
			WHERE type='table' AND name=?
		)`,
		tableName,
	).Scan(&exists)
	require.NoError(t, err)
	return exists
}

// IndexExists checks if an index exists in the database.
func IndexExists(t *testing.T, db *sql.DB, indexName string) bool {
	var exists bool
	err := db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM sqlite_master 
			WHERE type='index' AND name=?
		)`,
		indexName,
	).Scan(&exists)
	require.NoError(t, err)
	return exists
}

// ColumnExists checks if a column exists in a table.
func ColumnExists(t *testing.T, db *sql.DB, tableName, columnName string) bool {
	rows, err := db.Query(
		`SELECT 1 FROM pragma_table_info(?) WHERE name=?`,
		tableName, columnName,
	)
	require.NoError(t, err)
	defer rows.Close()

	return rows.Next()
}

// GetTableColumns returns a list of column names for a given table.
func GetTableColumns(t *testing.T, db *sql.DB, tableName string) []string {
	rows, err := db.Query(
		`SELECT name FROM pragma_table_info(?) ORDER BY cid`,
		tableName,
	)
	require.NoError(t, err)
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		require.NoError(t, err)
		columns = append(columns, name)
	}

	require.NoError(t, rows.Err())
	return columns
}

// InsertUser is a helper to insert a test user.
func InsertUser(t *testing.T, db *sql.DB, email, passwordHash, role string) int64 {
	result, err := db.Exec(
		`INSERT INTO users (email, password_hash, role) VALUES (?, ?, ?)`,
		email, passwordHash, role,
	)
	require.NoError(t, err)
	id, err := result.LastInsertId()
	require.NoError(t, err)
	return id
}

// InsertLabMember is a helper to insert a test lab member.
func InsertLabMember(t *testing.T, db *sql.DB, name, role string) int64 {
	result, err := db.Exec(
		`INSERT INTO lab_members (name, role) VALUES (?, ?)`,
		name, role,
	)
	require.NoError(t, err)
	id, err := result.LastInsertId()
	require.NoError(t, err)
	return id
}

// InsertPublication is a helper to insert a test publication.
func InsertPublication(t *testing.T, db *sql.DB, title, authorsText string, year int) int64 {
	result, err := db.Exec(
		`INSERT INTO publications (title, authors_text, year) VALUES (?, ?, ?)`,
		title, authorsText, year,
	)
	require.NoError(t, err)
	id, err := result.LastInsertId()
	require.NoError(t, err)
	return id
}

// InsertProject is a helper to insert a test project.
func InsertProject(t *testing.T, db *sql.DB, title, description, status string) int64 {
	result, err := db.Exec(
		`INSERT INTO projects (title, description, status) VALUES (?, ?, ?)`,
		title, description, status,
	)
	require.NoError(t, err)
	id, err := result.LastInsertId()
	require.NoError(t, err)
	return id
}
