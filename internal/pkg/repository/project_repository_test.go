package repository

import (
	"testing"

	"github.com/nekoteoj/lab-cms/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectRepository_CRUD(t *testing.T) {
	dbManager := setupTestDB(t)
	repo := NewProjectRepository(dbManager)

	t.Run("create project", func(t *testing.T) {
		proj := &models.Project{
			Title:       "Test Project",
			Description: "A test research project",
			Status:      models.ProjectStatusActive,
		}

		created, err := repo.Create(ctx, proj)
		require.NoError(t, err)
		assert.Greater(t, created.ID, 0)
		assert.Equal(t, "Test Project", created.Title)
		assert.Equal(t, models.ProjectStatusActive, created.Status)
	})

	t.Run("get project by id", func(t *testing.T) {
		proj := &models.Project{
			Title:       "Another Project",
			Description: "Description",
			Status:      models.ProjectStatusCompleted,
		}

		created, err := repo.Create(ctx, proj)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, "Another Project", retrieved.Title)
	})

	t.Run("get all projects", func(t *testing.T) {
		// Create projects with different statuses
		statuses := []models.ProjectStatus{
			models.ProjectStatusActive,
			models.ProjectStatusCompleted,
			models.ProjectStatusActive,
		}

		for i, status := range statuses {
			proj := &models.Project{
				Title:       "Project " + string(rune('A'+i)),
				Description: "Desc",
				Status:      status,
			}
			_, err := repo.Create(ctx, proj)
			require.NoError(t, err)
		}

		projects, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, projects)
	})

	t.Run("get by status", func(t *testing.T) {
		proj := &models.Project{
			Title:       "Active Project",
			Description: "Active",
			Status:      models.ProjectStatusActive,
		}

		_, err := repo.Create(ctx, proj)
		require.NoError(t, err)

		activeProjects, err := repo.GetByStatus(ctx, models.ProjectStatusActive)
		require.NoError(t, err)
		assert.NotEmpty(t, activeProjects)
	})

	t.Run("update project", func(t *testing.T) {
		proj := &models.Project{
			Title:       "Original Project",
			Description: "Original Desc",
			Status:      models.ProjectStatusActive,
		}

		created, err := repo.Create(ctx, proj)
		require.NoError(t, err)

		created.Title = "Updated Project"
		created.Status = models.ProjectStatusCompleted

		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)
		assert.Equal(t, "Updated Project", updated.Title)
		assert.Equal(t, models.ProjectStatusCompleted, updated.Status)
	})

	t.Run("delete project", func(t *testing.T) {
		proj := &models.Project{
			Title:       "To Delete",
			Description: "Will be deleted",
			Status:      models.ProjectStatusActive,
		}

		created, err := repo.Create(ctx, proj)
		require.NoError(t, err)

		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, created.ID)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestProjectRepository_Links(t *testing.T) {
	dbManager := setupTestDB(t)
	projRepo := NewProjectRepository(dbManager)
	memberRepo := NewLabMemberRepository(dbManager)
	pubRepo := NewPublicationRepository(dbManager)

	t.Run("link and unlink member", func(t *testing.T) {
		// Create project
		proj := &models.Project{
			Title:       "Team Project",
			Description: "Has members",
			Status:      models.ProjectStatusActive,
		}
		createdProj, err := projRepo.Create(ctx, proj)
		require.NoError(t, err)

		// Create member
		member := &models.LabMember{
			Name:         "Team Member",
			Role:         models.LabMemberRolePhD,
			DisplayOrder: 1,
		}
		createdMember, err := memberRepo.Create(ctx, member)
		require.NoError(t, err)

		// Link member
		err = projRepo.LinkMember(ctx, createdProj.ID, createdMember.ID)
		require.NoError(t, err)

		// Get members
		members, err := projRepo.GetMembers(ctx, createdProj.ID)
		require.NoError(t, err)
		assert.Len(t, members, 1)

		// Unlink member
		err = projRepo.UnlinkMember(ctx, createdProj.ID, createdMember.ID)
		require.NoError(t, err)

		// Verify unlink
		members, err = projRepo.GetMembers(ctx, createdProj.ID)
		require.NoError(t, err)
		assert.Empty(t, members)
	})

	t.Run("link and unlink publication", func(t *testing.T) {
		// Create project
		proj := &models.Project{
			Title:       "Research Project",
			Description: "Has publications",
			Status:      models.ProjectStatusActive,
		}
		createdProj, err := projRepo.Create(ctx, proj)
		require.NoError(t, err)

		// Create publication
		pub := &models.Publication{
			Title:       "Related Paper",
			AuthorsText: "Authors",
			Year:        2024,
		}
		createdPub, err := pubRepo.Create(ctx, pub)
		require.NoError(t, err)

		// Link publication
		err = projRepo.LinkPublication(ctx, createdProj.ID, createdPub.ID)
		require.NoError(t, err)

		// Get publications
		pubs, err := projRepo.GetPublications(ctx, createdProj.ID)
		require.NoError(t, err)
		assert.Len(t, pubs, 1)

		// Unlink publication
		err = projRepo.UnlinkPublication(ctx, createdProj.ID, createdPub.ID)
		require.NoError(t, err)

		// Verify unlink
		pubs, err = projRepo.GetPublications(ctx, createdProj.ID)
		require.NoError(t, err)
		assert.Empty(t, pubs)
	})

	t.Run("get project with relations", func(t *testing.T) {
		// Create project
		proj := &models.Project{
			Title:       "Full Project",
			Description: "Complete",
			Status:      models.ProjectStatusActive,
		}
		createdProj, err := projRepo.Create(ctx, proj)
		require.NoError(t, err)

		// Create member
		member := &models.LabMember{
			Name:         "Project Member",
			Role:         models.LabMemberRolePI,
			DisplayOrder: 2,
		}
		createdMember, err := memberRepo.Create(ctx, member)
		require.NoError(t, err)

		// Create publication
		pub := &models.Publication{
			Title:       "Project Paper",
			AuthorsText: "Authors",
			Year:        2023,
		}
		createdPub, err := pubRepo.Create(ctx, pub)
		require.NoError(t, err)

		// Link both
		err = projRepo.LinkMember(ctx, createdProj.ID, createdMember.ID)
		require.NoError(t, err)
		err = projRepo.LinkPublication(ctx, createdProj.ID, createdPub.ID)
		require.NoError(t, err)

		// Get with relations
		projWithRels, err := projRepo.GetWithRelations(ctx, createdProj.ID)
		require.NoError(t, err)
		assert.Equal(t, createdProj.ID, projWithRels.ID)
		assert.Len(t, projWithRels.Members, 1)
		assert.Len(t, projWithRels.Publications, 1)
	})
}
