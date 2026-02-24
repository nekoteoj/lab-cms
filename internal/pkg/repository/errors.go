package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"modernc.org/sqlite"
)

// SQLite extended error codes
// See: https://www.sqlite.org/rescode.html
const (
	sqliteConstraintUnique     = 2067 // SQLITE_CONSTRAINT_UNIQUE
	sqliteConstraintForeignKey = 787  // SQLITE_CONSTRAINT_FOREIGNKEY
	sqliteConstraintNotNull    = 1299 // SQLITE_CONSTRAINT_NOTNULL
	sqliteConstraintPrimaryKey = 1555 // SQLITE_CONSTRAINT_PRIMARYKEY
)

// Common errors that can be returned by repositories.
var (
	// ErrNotFound is returned when an entity is not found in the database.
	ErrNotFound = errors.New("entity not found")

	// ErrDuplicate is returned when attempting to create an entity that already exists.
	ErrDuplicate = errors.New("entity already exists")

	// ErrInvalidInput is returned when the input data is invalid.
	ErrInvalidInput = errors.New("invalid input")

	// ErrDatabase is returned for general database errors.
	ErrDatabase = errors.New("database error")
)

// isConstraintViolation checks if error is a specific SQLite constraint violation
func isConstraintViolation(err error, code int) bool {
	if err == nil {
		return false
	}

	// Try to get as sqlite.Error
	if sqliteErr, ok := err.(*sqlite.Error); ok {
		return sqliteErr.Code() == code
	}

	return false
}

// IsDuplicateError returns true if error is a unique or primary key constraint violation
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	return isConstraintViolation(err, sqliteConstraintUnique) ||
		isConstraintViolation(err, sqliteConstraintPrimaryKey)
}

// IsForeignKeyError returns true if error is a foreign key constraint violation
func IsForeignKeyError(err error) bool {
	return isConstraintViolation(err, sqliteConstraintForeignKey)
}

// IsNotNullError returns true if error is a NOT NULL constraint violation
func IsNotNullError(err error) bool {
	return isConstraintViolation(err, sqliteConstraintNotNull)
}

// WrapError wraps an error with operation context
// Returns specific error types for known error cases
func WrapError(err error, operation string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	// Check for duplicate errors
	if IsDuplicateError(err) {
		return ErrDuplicate
	}

	return fmt.Errorf("%s failed: %w", operation, err)
}
