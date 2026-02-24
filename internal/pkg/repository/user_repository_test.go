package repository

import (
	"testing"

	"github.com/nekoteoj/lab-cms/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_BasicOperations(t *testing.T) {
	dbManager := setupTestDB(t)
	repo := NewUserRepository(dbManager)

	t.Run("create user", func(t *testing.T) {
		user := &models.UserWithPassword{
			User: models.User{
				Email: "test@example.com",
				Role:  "normal",
			},
			PasswordHash: "hashedpassword123",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.Greater(t, created.ID, 0)
		assert.Equal(t, "test@example.com", created.Email)
		assert.Equal(t, models.UserRoleNormal, created.Role)
		assert.NotZero(t, created.CreatedAt)
		assert.NotZero(t, created.UpdatedAt)
	})

	t.Run("get user by id", func(t *testing.T) {
		user := &models.UserWithPassword{
			User: models.User{
				Email: "gettest@example.com",
				Role:  "normal",
			},
			PasswordHash: "hash",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Email, retrieved.Email)
	})

	t.Run("get user by email", func(t *testing.T) {
		user := &models.UserWithPassword{
			User: models.User{
				Email: "emailtest@example.com",
				Role:  "root",
			},
			PasswordHash: "hash",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		retrieved, err := repo.GetByEmail(ctx, "emailtest@example.com")
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Email, retrieved.Email)
		assert.Equal(t, "hash", retrieved.PasswordHash)
	})

	t.Run("get all users", func(t *testing.T) {
		// Clear existing users first
		db := dbManager.GetDB()
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)

		// Create multiple users
		for i := 0; i < 3; i++ {
			user := &models.UserWithPassword{
				User: models.User{
					Email: string('a'+byte(i)) + "@example.com",
					Role:  "normal",
				},
				PasswordHash: "hash",
			}
			_, err := repo.Create(ctx, user)
			require.NoError(t, err)
		}

		users, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, users, 3)
	})

	t.Run("update user", func(t *testing.T) {
		user := &models.UserWithPassword{
			User: models.User{
				Email: "update@example.com",
				Role:  "normal",
			},
			PasswordHash: "hash",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		created.Email = "updated@example.com"
		created.Role = models.UserRoleRoot

		updated, err := repo.Update(ctx, &created.User)
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", updated.Email)
		assert.Equal(t, models.UserRoleRoot, updated.Role)
	})

	t.Run("delete user", func(t *testing.T) {
		user := &models.UserWithPassword{
			User: models.User{
				Email: "delete@example.com",
				Role:  "normal",
			},
			PasswordHash: "hash",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, created.ID)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("duplicate email error", func(t *testing.T) {
		user1 := &models.UserWithPassword{
			User: models.User{
				Email: "duplicate@example.com",
				Role:  "normal",
			},
			PasswordHash: "hash1",
		}

		_, err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := &models.UserWithPassword{
			User: models.User{
				Email: "duplicate@example.com",
				Role:  "normal",
			},
			PasswordHash: "hash2",
		}

		_, err = repo.Create(ctx, user2)
		assert.Equal(t, ErrDuplicate, err)
	})

	t.Run("not found errors", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 99999)
		assert.Equal(t, ErrNotFound, err)

		_, err = repo.GetByEmail(ctx, "nonexistent@example.com")
		assert.Equal(t, ErrNotFound, err)

		err = repo.Delete(ctx, 99999)
		assert.Equal(t, ErrNotFound, err)
	})
}
