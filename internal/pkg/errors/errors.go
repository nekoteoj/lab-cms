// Package errors provides application-specific error types and error handling utilities.
// These error types map to appropriate HTTP status codes and provide context for debugging.
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Application error types
var (
	// ErrNotFound is returned when a requested resource doesn't exist
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput is returned when user input fails validation
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized is returned when authentication is required but missing or invalid
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when the user lacks permission for an action
	ErrForbidden = errors.New("forbidden")

	// ErrInternal is returned for unexpected server errors
	ErrInternal = errors.New("internal server error")

	// ErrDuplicate is returned when attempting to create a resource that already exists
	ErrDuplicate = errors.New("resource already exists")

	// ErrDatabase is returned for database-related errors
	ErrDatabase = errors.New("database error")
)

// AppError is a custom error type that includes an HTTP status code
// and additional context for debugging
type AppError struct {
	// Code is a machine-readable error code
	Code string `json:"code"`

	// Message is a human-readable error message (safe to show to users)
	Message string `json:"message"`

	// StatusCode is the HTTP status code to return
	StatusCode int `json:"-"`

	// Details contains additional context for debugging (not exposed to users in production)
	Details string `json:"details,omitempty"`

	// Cause is the underlying error that caused this error
	Cause error `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error for error chain inspection
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError creates a new AppError with the given parameters
func NewAppError(code string, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// Wrap wraps an existing error with context
func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Cause:      err,
	}
}

// WithDetails adds debugging details to the error
func (e *AppError) WithDetails(details string) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Details:    details,
		Cause:      e.Cause,
	}
}

// Predefined error constructors for common cases

// NotFound creates a not found error for the given resource
func NotFound(resource string, identifier interface{}) *AppError {
	var details string
	if identifier == nil {
		details = fmt.Sprintf("%s does not exist", resource)
	} else {
		details = fmt.Sprintf("%s with identifier '%v' does not exist", resource, identifier)
	}

	return &AppError{
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
		Details:    details,
	}
}

// Validation creates a validation error with a specific message
func Validation(field string, issue string) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    fmt.Sprintf("Invalid %s: %s", field, issue),
		StatusCode: http.StatusBadRequest,
		Details:    fmt.Sprintf("Field '%s' failed validation: %s", field, issue),
	}
}

// ValidationFromErr creates a validation error from an existing error
func ValidationFromErr(err error) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    "Validation failed",
		StatusCode: http.StatusBadRequest,
		Cause:      err,
	}
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	if message == "" {
		message = "Authentication required"
	}
	return &AppError{
		Code:       "UNAUTHORIZED",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// Forbidden creates a forbidden error
func Forbidden(action string) *AppError {
	return &AppError{
		Code:       "FORBIDDEN",
		Message:    "You don't have permission to perform this action",
		StatusCode: http.StatusForbidden,
		Details:    fmt.Sprintf("Action '%s' is not permitted", action),
	}
}

// Internal creates an internal server error
func Internal(err error) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "An unexpected error occurred. Please try again later.",
		StatusCode: http.StatusInternalServerError,
		Cause:      err,
	}
}

// Duplicate creates a duplicate resource error
func Duplicate(resource string, field string) *AppError {
	return &AppError{
		Code:       "DUPLICATE_ERROR",
		Message:    fmt.Sprintf("A %s with this %s already exists", resource, field),
		StatusCode: http.StatusConflict,
		Details:    fmt.Sprintf("Duplicate %s on field '%s'", resource, field),
	}
}

// Database creates a database error
func Database(err error) *AppError {
	return &AppError{
		Code:       "DATABASE_ERROR",
		Message:    "A database error occurred. Please try again later.",
		StatusCode: http.StatusInternalServerError,
		Cause:      err,
	}
}

// Error checking helpers

// IsNotFound returns true if the error is a not found error
func IsNotFound(err error) bool {
	return isError(err, ErrNotFound) || hasStatusCode(err, http.StatusNotFound)
}

// IsValidationError returns true if the error is a validation error
func IsValidationError(err error) bool {
	return isError(err, ErrInvalidInput) || hasStatusCode(err, http.StatusBadRequest)
}

// IsUnauthorized returns true if the error is an unauthorized error
func IsUnauthorized(err error) bool {
	return isError(err, ErrUnauthorized) || hasStatusCode(err, http.StatusUnauthorized)
}

// IsForbidden returns true if the error is a forbidden error
func IsForbidden(err error) bool {
	return isError(err, ErrForbidden) || hasStatusCode(err, http.StatusForbidden)
}

// IsInternalError returns true if the error is an internal error
func IsInternalError(err error) bool {
	return isError(err, ErrInternal) || hasStatusCode(err, http.StatusInternalServerError)
}

// IsDuplicate returns true if the error is a duplicate resource error
func IsDuplicate(err error) bool {
	return isError(err, ErrDuplicate) || hasStatusCode(err, http.StatusConflict)
}

// isError checks if the error matches the target error (unwraps the chain)
func isError(err, target error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, target)
}

// hasStatusCode checks if the error is an AppError with the given status code
func hasStatusCode(err error, code int) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode == code
	}
	return false
}

// GetStatusCode extracts the HTTP status code from an error
// Returns 500 if the error is not an AppError
func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}

	return http.StatusInternalServerError
}

// Wrap adds context to an error using fmt.Errorf with %w
func Wrap(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}
