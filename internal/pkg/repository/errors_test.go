package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsDuplicateError(t *testing.T) {
	dbManager := setupTestDB(t)
	db := dbManager.GetDB()

	t.Run("returns false for nil error", func(t *testing.T) {
		result := IsDuplicateError(nil)
		assert.False(t, result)
	})

	t.Run("returns false for non-SQLite error", func(t *testing.T) {
		err := errors.New("some random error")
		result := IsDuplicateError(err)
		assert.False(t, result)
	})

	t.Run("detects actual unique constraint violation from database", func(t *testing.T) {
		// Create a table with unique constraint
		_, err := db.Exec(`CREATE TABLE test_unique (id INTEGER PRIMARY KEY, email TEXT UNIQUE)`)
		require.NoError(t, err)

		// Insert first record
		_, err = db.Exec(`INSERT INTO test_unique (email) VALUES ('test@example.com')`)
		require.NoError(t, err)

		// Try to insert duplicate - should get unique constraint error
		_, err = db.Exec(`INSERT INTO test_unique (email) VALUES ('test@example.com')`)
		require.Error(t, err)

		result := IsDuplicateError(err)
		assert.True(t, result, "Should detect unique constraint violation, got error: %v", err)
	})

	t.Run("detects primary key constraint violation", func(t *testing.T) {
		// Create table with explicit primary key
		_, err := db.Exec(`CREATE TABLE test_pk (id INTEGER PRIMARY KEY, name TEXT)`)
		require.NoError(t, err)

		// Insert first record
		_, err = db.Exec(`INSERT INTO test_pk (id, name) VALUES (1, 'first')`)
		require.NoError(t, err)

		// Try to insert duplicate primary key
		_, err = db.Exec(`INSERT INTO test_pk (id, name) VALUES (1, 'second')`)
		require.Error(t, err)

		result := IsDuplicateError(err)
		assert.True(t, result, "Should detect primary key constraint violation")
	})
}

func TestIsForeignKeyError(t *testing.T) {
	dbManager := setupTestDB(t)
	db := dbManager.GetDB()

	t.Run("returns false for nil error", func(t *testing.T) {
		result := IsForeignKeyError(nil)
		assert.False(t, result)
	})

	t.Run("detects actual foreign key constraint violation", func(t *testing.T) {
		// Create parent table
		_, err := db.Exec(`CREATE TABLE parent (id INTEGER PRIMARY KEY)`)
		require.NoError(t, err)

		// Create child table with foreign key
		_, err = db.Exec(`CREATE TABLE child (id INTEGER PRIMARY KEY, parent_id INTEGER REFERENCES parent(id))`)
		require.NoError(t, err)

		// Try to insert child with non-existent parent
		_, err = db.Exec(`INSERT INTO child (parent_id) VALUES (999)`)
		require.Error(t, err)

		result := IsForeignKeyError(err)
		assert.True(t, result, "Should detect foreign key constraint violation")
	})

	t.Run("returns false for other errors", func(t *testing.T) {
		err := errors.New("random error")
		result := IsForeignKeyError(err)
		assert.False(t, result)
	})
}

func TestIsNotNullError(t *testing.T) {
	dbManager := setupTestDB(t)
	db := dbManager.GetDB()

	t.Run("returns false for nil error", func(t *testing.T) {
		result := IsNotNullError(nil)
		assert.False(t, result)
	})

	t.Run("detects NOT NULL constraint violation", func(t *testing.T) {
		// Create table with NOT NULL constraint
		_, err := db.Exec(`CREATE TABLE test_notnull (id INTEGER PRIMARY KEY, name TEXT NOT NULL)`)
		require.NoError(t, err)

		// Try to insert NULL value
		_, err = db.Exec(`INSERT INTO test_notnull (name) VALUES (NULL)`)
		require.Error(t, err)

		result := IsNotNullError(err)
		assert.True(t, result, "Should detect NOT NULL constraint violation")
	})

	t.Run("returns false for other errors", func(t *testing.T) {
		err := errors.New("random error")
		result := IsNotNullError(err)
		assert.False(t, result)
	})
}

func TestWrapError(t *testing.T) {
	t.Run("returns nil for nil error", func(t *testing.T) {
		result := WrapError(nil, "operation")
		assert.Nil(t, result)
	})

	t.Run("returns ErrNotFound for sql.ErrNoRows", func(t *testing.T) {
		err := WrapError(sql.ErrNoRows, "get user")
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns ErrDuplicate for unique constraint", func(t *testing.T) {
		// Create a mock SQLite error by actually causing one
		dbManager := setupTestDB(t)
		db := dbManager.GetDB()

		_, err := db.Exec(`CREATE TABLE test_wrap (email TEXT UNIQUE)`)
		require.NoError(t, err)

		_, err = db.Exec(`INSERT INTO test_wrap VALUES ('dup')`)
		require.NoError(t, err)

		_, err = db.Exec(`INSERT INTO test_wrap VALUES ('dup')`)
		require.Error(t, err)

		wrapped := WrapError(err, "create user")
		assert.Equal(t, ErrDuplicate, wrapped)
	})

	t.Run("wraps unknown errors with context", func(t *testing.T) {
		originalErr := errors.New("connection refused")
		err := WrapError(originalErr, "database query")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database query failed")
		assert.Contains(t, err.Error(), "connection refused")
	})
}

func TestErrorConstants(t *testing.T) {
	t.Run("ErrNotFound is defined", func(t *testing.T) {
		assert.NotNil(t, ErrNotFound)
		assert.Equal(t, "entity not found", ErrNotFound.Error())
	})

	t.Run("ErrDuplicate is defined", func(t *testing.T) {
		assert.NotNil(t, ErrDuplicate)
		assert.Equal(t, "entity already exists", ErrDuplicate.Error())
	})

	t.Run("ErrInvalidInput is defined", func(t *testing.T) {
		assert.NotNil(t, ErrInvalidInput)
		assert.Equal(t, "invalid input", ErrInvalidInput.Error())
	})

	t.Run("ErrDatabase is defined", func(t *testing.T) {
		assert.NotNil(t, ErrDatabase)
		assert.Equal(t, "database error", ErrDatabase.Error())
	})
}
