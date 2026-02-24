package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appErr   *AppError
		expected string
	}{
		{
			name:     "with cause",
			appErr:   NewAppError("TEST", "test message", http.StatusBadRequest).Wrap(errors.New("cause")),
			expected: "test message: cause",
		},
		{
			name:     "without cause",
			appErr:   NewAppError("TEST", "test message", http.StatusBadRequest),
			expected: "test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.appErr.Error(); got != tt.expected {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("original error")
	appErr := NewAppError("TEST", "wrapped", http.StatusBadRequest).Wrap(cause)

	if unwrapped := appErr.Unwrap(); unwrapped != cause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestNewAppError(t *testing.T) {
	appErr := NewAppError("CODE123", "test message", http.StatusBadRequest)

	if appErr.Code != "CODE123" {
		t.Errorf("Code = %v, want %v", appErr.Code, "CODE123")
	}
	if appErr.Message != "test message" {
		t.Errorf("Message = %v, want %v", appErr.Message, "test message")
	}
	if appErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %v, want %v", appErr.StatusCode, http.StatusBadRequest)
	}
}

func TestAppError_Wrap(t *testing.T) {
	original := NewAppError("ORIGINAL", "original", http.StatusBadRequest)
	cause := errors.New("cause error")

	wrapped := original.Wrap(cause)

	if wrapped.Code != original.Code {
		t.Error("Wrap should preserve Code")
	}
	if wrapped.Message != original.Message {
		t.Error("Wrap should preserve Message")
	}
	if wrapped.StatusCode != original.StatusCode {
		t.Error("Wrap should preserve StatusCode")
	}
	if wrapped.Cause != cause {
		t.Error("Wrap should set Cause")
	}
}

func TestAppError_WithDetails(t *testing.T) {
	appErr := NewAppError("CODE", "message", http.StatusBadRequest).
		WithDetails("debug details")

	if appErr.Details != "debug details" {
		t.Errorf("Details = %v, want %v", appErr.Details, "debug details")
	}
}

func TestNotFound(t *testing.T) {
	err := NotFound("User", 123)

	if err.Code != "NOT_FOUND" {
		t.Errorf("Code = %v, want NOT_FOUND", err.Code)
	}
	if err.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %v, want 404", err.StatusCode)
	}
	if err.Message != "User not found" {
		t.Errorf("Message = %v, want 'User not found'", err.Message)
	}
	if err.Details != "User with identifier '123' does not exist" {
		t.Errorf("Details = %v, want details about identifier", err.Details)
	}
}

func TestValidation(t *testing.T) {
	err := Validation("email", "invalid format")

	if err.Code != "VALIDATION_ERROR" {
		t.Errorf("Code = %v, want VALIDATION_ERROR", err.Code)
	}
	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %v, want 400", err.StatusCode)
	}
	if err.Message != "Invalid email: invalid format" {
		t.Errorf("Message = %v", err.Message)
	}
}

func TestValidationFromErr(t *testing.T) {
	originalErr := errors.New("parse error")
	err := ValidationFromErr(originalErr)

	if err.Code != "VALIDATION_ERROR" {
		t.Errorf("Code = %v, want VALIDATION_ERROR", err.Code)
	}
	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %v, want 400", err.StatusCode)
	}
	if err.Cause != originalErr {
		t.Error("Should wrap original error")
	}
}

func TestUnauthorized(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{"with message", "custom message", "custom message"},
		{"without message", "", "Authentication required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unauthorized(tt.message)
			if err.Code != "UNAUTHORIZED" {
				t.Errorf("Code = %v, want UNAUTHORIZED", err.Code)
			}
			if err.StatusCode != http.StatusUnauthorized {
				t.Errorf("StatusCode = %v, want 401", err.StatusCode)
			}
			if err.Message != tt.expected {
				t.Errorf("Message = %v, want %v", err.Message, tt.expected)
			}
		})
	}
}

func TestForbidden(t *testing.T) {
	err := Forbidden("delete_user")

	if err.Code != "FORBIDDEN" {
		t.Errorf("Code = %v, want FORBIDDEN", err.Code)
	}
	if err.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %v, want 403", err.StatusCode)
	}
	if err.Details != "Action 'delete_user' is not permitted" {
		t.Errorf("Details = %v", err.Details)
	}
}

func TestInternal(t *testing.T) {
	cause := errors.New("database connection failed")
	err := Internal(cause)

	if err.Code != "INTERNAL_ERROR" {
		t.Errorf("Code = %v, want INTERNAL_ERROR", err.Code)
	}
	if err.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %v, want 500", err.StatusCode)
	}
	if err.Cause != cause {
		t.Error("Should wrap original error")
	}
}

func TestDuplicate(t *testing.T) {
	err := Duplicate("User", "email")

	if err.Code != "DUPLICATE_ERROR" {
		t.Errorf("Code = %v, want DUPLICATE_ERROR", err.Code)
	}
	if err.StatusCode != http.StatusConflict {
		t.Errorf("StatusCode = %v, want 409", err.StatusCode)
	}
	if err.Message != "A User with this email already exists" {
		t.Errorf("Message = %v", err.Message)
	}
}

func TestDatabase(t *testing.T) {
	cause := errors.New("connection timeout")
	err := Database(cause)

	if err.Code != "DATABASE_ERROR" {
		t.Errorf("Code = %v, want DATABASE_ERROR", err.Code)
	}
	if err.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %v, want 500", err.StatusCode)
	}
	if err.Cause != cause {
		t.Error("Should wrap original error")
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected bool
	}{
		{"nil error", nil, false},
		{"not found error", ErrNotFound, true},
		{"not found app error", NotFound("User", 1), true},
		{"other error", errors.New("other"), false},
		{"regular 404", errors.New("404"), false}, // Not an AppError
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.input); got != tt.expected {
				t.Errorf("IsNotFound(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected bool
	}{
		{"nil error", nil, false},
		{"validation error", ErrInvalidInput, true},
		{"validation app error", Validation("field", "error"), true},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidationError(tt.input); got != tt.expected {
				t.Errorf("IsValidationError(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected bool
	}{
		{"nil error", nil, false},
		{"unauthorized error", ErrUnauthorized, true},
		{"unauthorized app error", Unauthorized(""), true},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnauthorized(tt.input); got != tt.expected {
				t.Errorf("IsUnauthorized(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsForbidden(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected bool
	}{
		{"nil error", nil, false},
		{"forbidden error", ErrForbidden, true},
		{"forbidden app error", Forbidden("action"), true},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsForbidden(tt.input); got != tt.expected {
				t.Errorf("IsForbidden(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsInternalError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected bool
	}{
		{"nil error", nil, false},
		{"internal error", ErrInternal, true},
		{"internal app error", Internal(errors.New("cause")), true},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInternalError(tt.input); got != tt.expected {
				t.Errorf("IsInternalError(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected bool
	}{
		{"nil error", nil, false},
		{"duplicate error", ErrDuplicate, true},
		{"duplicate app error", Duplicate("User", "email"), true},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDuplicate(tt.input); got != tt.expected {
				t.Errorf("IsDuplicate(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGetStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected int
	}{
		{"nil error", nil, http.StatusOK},
		{"not found", NotFound("User", 1), http.StatusNotFound},
		{"validation error", Validation("field", "error"), http.StatusBadRequest},
		{"unauthorized", Unauthorized(""), http.StatusUnauthorized},
		{"forbidden", Forbidden("action"), http.StatusForbidden},
		{"internal error", Internal(errors.New("cause")), http.StatusInternalServerError},
		{"duplicate", Duplicate("User", "email"), http.StatusConflict},
		{"regular error", errors.New("other"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStatusCode(tt.input); got != tt.expected {
				t.Errorf("GetStatusCode(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name          string
		input         error
		format        string
		args          []interface{}
		expectedNil   bool
		expectedError string
	}{
		{
			name:        "nil error returns nil",
			input:       nil,
			format:      "context",
			expectedNil: true,
		},
		{
			name:          "wrap with format",
			input:         errors.New("original"),
			format:        "failed to do %s",
			args:          []interface{}{"something"},
			expectedError: "failed to do something: original",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.input, tt.format, tt.args...)
			if tt.expectedNil {
				if got != nil {
					t.Errorf("Wrap() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("Wrap() returned nil, expected error")
			}
			if got.Error() != tt.expectedError {
				t.Errorf("Wrap() = %v, want %v", got.Error(), tt.expectedError)
			}
			// Test that we can unwrap
			if !errors.Is(got, tt.input) {
				t.Error("Wrapped error should be unwrappable")
			}
		})
	}
}

func TestErrorsIs(t *testing.T) {
	// Test that standard errors.Is works with our error types
	appErr := NewAppError("CODE", "message", http.StatusBadRequest).Wrap(ErrNotFound)

	if !errors.Is(appErr, ErrNotFound) {
		t.Error("errors.Is should find wrapped error")
	}

	if !errors.Is(appErr, appErr) {
		t.Error("errors.Is should match the error itself")
	}
}

func TestErrorsAs(t *testing.T) {
	// Test that errors.As works with our error types
	appErr := NotFound("User", 123)
	var target *AppError

	if !errors.As(appErr, &target) {
		t.Error("errors.As should match *AppError")
	}

	if target.Code != "NOT_FOUND" {
		t.Errorf("Target.Code = %v, want NOT_FOUND", target.Code)
	}
}

func TestPredefinedErrors(t *testing.T) {
	// Test that predefined errors are usable
	errors := []error{
		ErrNotFound,
		ErrInvalidInput,
		ErrUnauthorized,
		ErrForbidden,
		ErrInternal,
		ErrDuplicate,
		ErrDatabase,
	}

	for _, err := range errors {
		if err == nil {
			t.Error("Predefined error should not be nil")
		}
	}
}
