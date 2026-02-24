package repository

import (
	"database/sql"
	"testing"

	"github.com/nekoteoj/lab-cms/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabMemberRepository_CRUD(t *testing.T) {
	dbManager := setupTestDB(t)
	repo := NewLabMemberRepository(dbManager)

	t.Run("create member", func(t *testing.T) {
		member := &models.LabMember{
			Name:         "Dr. John Doe",
			Role:         models.LabMemberRolePI,
			Email:        sql.NullString{String: "john@example.com", Valid: true},
			Bio:          sql.NullString{String: "Principal Investigator", Valid: true},
			DisplayOrder: 1,
		}

		created, err := repo.Create(ctx, member)
		require.NoError(t, err)
		assert.Greater(t, created.ID, 0)
		assert.Equal(t, "Dr. John Doe", created.Name)
		assert.Equal(t, models.LabMemberRolePI, created.Role)
	})

	t.Run("get member by id", func(t *testing.T) {
		member := &models.LabMember{
			Name:         "Jane Smith",
			Role:         models.LabMemberRolePhD,
			DisplayOrder: 2,
		}

		created, err := repo.Create(ctx, member)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, "Jane Smith", retrieved.Name)
	})

	t.Run("get all members", func(t *testing.T) {
		// Create multiple members
		roles := []models.LabMemberRole{
			models.LabMemberRolePI,
			models.LabMemberRolePhD,
			models.LabMemberRolePostdoc,
		}

		for i, role := range roles {
			member := &models.LabMember{
				Name:         string('A'+byte(i)) + " Member",
				Role:         role,
				DisplayOrder: i,
			}
			_, err := repo.Create(ctx, member)
			require.NoError(t, err)
		}

		members, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, members)
	})

	t.Run("get by role", func(t *testing.T) {
		member := &models.LabMember{
			Name:         "Researcher Test",
			Role:         models.LabMemberRoleResearcher,
			DisplayOrder: 10,
		}

		_, err := repo.Create(ctx, member)
		require.NoError(t, err)

		members, err := repo.GetByRole(ctx, models.LabMemberRoleResearcher)
		require.NoError(t, err)
		assert.NotEmpty(t, members)
	})

	t.Run("update member", func(t *testing.T) {
		member := &models.LabMember{
			Name:         "Original Name",
			Role:         models.LabMemberRoleMaster,
			DisplayOrder: 5,
		}

		created, err := repo.Create(ctx, member)
		require.NoError(t, err)

		created.Name = "Updated Name"
		created.Bio = sql.NullString{String: "Updated bio", Valid: true}

		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)
		assert.Equal(t, "Updated bio", updated.Bio.String)
	})

	t.Run("mark as alumni", func(t *testing.T) {
		member := &models.LabMember{
			Name:         "Former Member",
			Role:         models.LabMemberRolePhD,
			IsAlumni:     false,
			DisplayOrder: 20,
		}

		created, err := repo.Create(ctx, member)
		require.NoError(t, err)

		err = repo.MarkAsAlumni(ctx, created.ID, true)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.True(t, retrieved.IsAlumni)
	})

	t.Run("get alumni", func(t *testing.T) {
		// Create an alumni member
		member := &models.LabMember{
			Name:         "Alumni Member",
			Role:         models.LabMemberRolePhD,
			IsAlumni:     true,
			DisplayOrder: 30,
		}

		_, err := repo.Create(ctx, member)
		require.NoError(t, err)

		alumni, err := repo.GetAlumni(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, alumni)
	})

	t.Run("delete member", func(t *testing.T) {
		member := &models.LabMember{
			Name:         "To Delete",
			Role:         models.LabMemberRoleBachelor,
			DisplayOrder: 40,
		}

		created, err := repo.Create(ctx, member)
		require.NoError(t, err)

		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, created.ID)
		assert.Equal(t, ErrNotFound, err)
	})
}
