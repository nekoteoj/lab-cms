package repository

import (
	"context"

	"github.com/nekoteoj/lab-cms/internal/pkg/models"
)

// Repository is the base interface for generic CRUD repositories
// T is the entity type (e.g., models.User, models.Project)
type Repository[T any] interface {
	// GetByID retrieves a single entity by its ID
	// Returns ErrNotFound if the entity does not exist
	GetByID(ctx context.Context, id int) (*T, error)

	// GetAll retrieves all entities
	GetAll(ctx context.Context) ([]T, error)

	// Create inserts a new entity and returns it with the generated ID
	Create(ctx context.Context, entity *T) (*T, error)

	// Update modifies an existing entity
	// Returns ErrNotFound if the entity does not exist
	Update(ctx context.Context, entity *T) (*T, error)

	// Delete removes an entity by its ID
	// Returns ErrNotFound if the entity does not exist
	Delete(ctx context.Context, id int) error
}

// UserAuthRepository is a specialized interface for user authentication
// This extends basic operations with authentication-specific methods
type UserAuthRepository interface {
	// Basic CRUD (returns User without password for security)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id int) error

	// Authentication-specific methods (handles password securely)
	GetByEmail(ctx context.Context, email string) (*models.UserWithPassword, error)
	Create(ctx context.Context, user *models.UserWithPassword) (*models.UserWithPassword, error)
	UpdatePassword(ctx context.Context, id int, passwordHash string) error
}
