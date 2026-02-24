package repository

import (
	"context"
	"database/sql"

	"github.com/nekoteoj/lab-cms/internal/pkg/db"
	"github.com/nekoteoj/lab-cms/internal/pkg/models"
)

// Ensure UserRepository implements UserAuthRepository interface
var _ UserAuthRepository = (*UserRepository)(nil)

// UserRepository provides data access for users.
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new user repository.
func NewUserRepository(dbManager *db.DBManager) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(dbManager, "users"),
	}
}

// GetByID retrieves a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.GetExecer(ctx).QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, WrapError(err, "get user by id")
	}

	return &user, nil
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.UserWithPassword, error) {
	query := `
		SELECT id, email, role, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := r.GetExecer(ctx).QueryRowContext(ctx, query, email)

	var user models.UserWithPassword
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, WrapError(err, "get user by email")
	}

	return &user, nil
}

// GetAll retrieves all users.
func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, email, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, WrapError(err, "get all users")
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan user")
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate users")
	}

	return users, nil
}

// Create inserts a new user.
func (r *UserRepository) Create(ctx context.Context, user *models.UserWithPassword) (*models.UserWithPassword, error) {
	query := `
		INSERT INTO users (email, role, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, datetime('now'), datetime('now'))
		RETURNING id, created_at, updated_at
	`

	row := r.GetExecer(ctx).QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Role,
		user.PasswordHash,
	)

	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if IsDuplicateError(err) {
			return nil, ErrDuplicate
		}
		return nil, WrapError(err, "create user")
	}

	return user, nil
}

// Update modifies an existing user.
func (r *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		UPDATE users
		SET email = $1, role = $2, updated_at = datetime('now')
		WHERE id = $3
		RETURNING updated_at
	`

	row := r.GetExecer(ctx).QueryRowContext(ctx, query, user.Email, user.Role, user.ID)

	err := row.Scan(&user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		if IsDuplicateError(err) {
			return nil, ErrDuplicate
		}
		return nil, WrapError(err, "update user")
	}

	return user, nil
}

// UpdatePassword updates a user's password.
func (r *UserRepository) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = datetime('now')
		WHERE id = $2
	`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, passwordHash, id)
	if err != nil {
		return WrapError(err, "update password")
	}

	return CheckRowsAffected(result, 1)
}

// Delete removes a user.
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return WrapError(err, "delete user")
	}

	return CheckRowsAffected(result, 1)
}
