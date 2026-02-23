package models

import (
	"time"
)

// User represents an admin user in the system
// Password hash is handled separately for security
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email" validate:"required,email,max=255"`
	Role      UserRole  `json:"role" validate:"required,oneof=normal root"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserWithPassword extends User to include password for authentication
// This should only be used in authentication contexts, not in general API responses
type UserWithPassword struct {
	User
	PasswordHash string `json:"-"` // Never serialized to JSON
}
