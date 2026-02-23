package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Validation(t *testing.T) {
	v := newValidator()

	validUser := User{
		Email: "test@example.com",
		Role:  UserRoleNormal,
	}

	err := validateStruct(v, validUser)
	assert.NoError(t, err, "valid user should pass validation")
}

func TestUser_Validation_InvalidEmail(t *testing.T) {
	v := newValidator()

	tests := []struct {
		name  string
		email string
	}{
		{"empty email", ""},
		{"missing @", "invalidemail"},
		{"missing domain", "test@"},
		{"invalid format", "test@.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{
				Email: tt.email,
				Role:  UserRoleNormal,
			}
			err := validateStruct(v, user)
			assert.Error(t, err, "invalid email should fail validation")
		})
	}
}

func TestUser_Validation_InvalidRole(t *testing.T) {
	v := newValidator()

	tests := []struct {
		name string
		role string
	}{
		{"empty role", ""},
		{"invalid role", "admin"},
		{"wrong case", "Normal"},
		{"typo", "rooot"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{
				Email: "test@example.com",
				Role:  UserRole(tt.role),
			}
			err := validateStruct(v, user)
			assert.Error(t, err, "invalid role should fail validation")
		})
	}
}

func TestUser_Validation_MissingRequired(t *testing.T) {
	v := newValidator()

	user := User{}
	err := validateStruct(v, user)
	assert.Error(t, err, "empty user should fail validation")
	assert.Contains(t, err.Error(), "Email")
	assert.Contains(t, err.Error(), "Role")
}

func TestUser_JSONSerialization(t *testing.T) {
	user := User{
		ID:    1,
		Email: "test@example.com",
		Role:  UserRoleNormal,
	}

	data, err := json.Marshal(user)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"id\":1")
	assert.Contains(t, jsonStr, "\"email\":\"test@example.com\"")
	assert.Contains(t, jsonStr, "\"role\":\"normal\"")
	assert.Contains(t, jsonStr, "\"created_at\"")
}

func TestUser_JSON_DoesNotIncludePassword(t *testing.T) {
	user := User{
		ID:    1,
		Email: "test@example.com",
		Role:  UserRoleRoot,
	}

	data, err := json.Marshal(user)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.NotContains(t, jsonStr, "password", "password should not appear in JSON")
}

func TestUserWithPassword(t *testing.T) {
	v := newValidator()

	user := UserWithPassword{
		User: User{
			ID:    1,
			Email: "admin@example.com",
			Role:  UserRoleRoot,
		},
		PasswordHash: "$2a$10$hashedpasswordhere",
	}

	// Validate embedded User fields
	err := validateStruct(v, user.User)
	assert.NoError(t, err)

	// JSON should not include password
	data, err := json.Marshal(user)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "admin@example.com")
	assert.NotContains(t, strings.ToLower(jsonStr), "password")
	assert.NotContains(t, jsonStr, "PasswordHash")
	assert.NotContains(t, jsonStr, "$2a$10")
}

func TestUser_JSONDeserialization(t *testing.T) {
	jsonData := `{"id":1,"email":"test@example.com","role":"root","created_at":"2024-01-01T00:00:00Z"}`

	var user User
	err := json.Unmarshal([]byte(jsonData), &user)
	require.NoError(t, err)

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, UserRoleRoot, user.Role)
}
