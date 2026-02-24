package repository

import (
	"database/sql"
	"testing"

	"github.com/nekoteoj/lab-cms/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicationRepository_CRUD(t *testing.T) {
	dbManager := setupTestDB(t)
	repo := NewPublicationRepository(dbManager)

	t.Run("create publication", func(t *testing.T) {
		pub := &models.Publication{
			Title:       "Test Publication",
			AuthorsText: "John Doe, Jane Smith",
			Venue:       sql.NullString{String: "Nature", Valid: true},
			Year:        2024,
			URL:         sql.NullString{String: "https://example.com/paper", Valid: true},
		}

		created, err := repo.Create(ctx, pub)
		require.NoError(t, err)
		assert.Greater(t, created.ID, 0)
		assert.Equal(t, "Test Publication", created.Title)
		assert.Equal(t, 2024, created.Year)
	})

	t.Run("get publication by id", func(t *testing.T) {
		pub := &models.Publication{
			Title:       "Another Paper",
			AuthorsText: "Alice Bob",
			Year:        2023,
		}

		created, err := repo.Create(ctx, pub)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, "Another Paper", retrieved.Title)
	})

	t.Run("get all publications", func(t *testing.T) {
		// Create publications in different years
		years := []int{2024, 2023, 2022}
		for _, year := range years {
			pub := &models.Publication{
				Title:       "Paper from " + string(rune('0'+year%10)),
				AuthorsText: "Author",
				Year:        year,
			}
			_, err := repo.Create(ctx, pub)
			require.NoError(t, err)
		}

		pubs, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, pubs)

		// Should be sorted by year descending
		if len(pubs) >= 2 {
			assert.GreaterOrEqual(t, pubs[0].Year, pubs[1].Year)
		}
	})

	t.Run("get by year", func(t *testing.T) {
		pub := &models.Publication{
			Title:       "2021 Paper",
			AuthorsText: "Author",
			Year:        2021,
		}

		_, err := repo.Create(ctx, pub)
		require.NoError(t, err)

		pubs, err := repo.GetByYear(ctx, 2021)
		require.NoError(t, err)
		assert.NotEmpty(t, pubs)
	})

	t.Run("update publication", func(t *testing.T) {
		pub := &models.Publication{
			Title:       "Original Title",
			AuthorsText: "Original Authors",
			Year:        2020,
		}

		created, err := repo.Create(ctx, pub)
		require.NoError(t, err)

		created.Title = "Updated Title"
		created.Year = 2021

		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, 2021, updated.Year)
	})

	t.Run("delete publication", func(t *testing.T) {
		pub := &models.Publication{
			Title:       "To Delete",
			AuthorsText: "Authors",
			Year:        2019,
		}

		created, err := repo.Create(ctx, pub)
		require.NoError(t, err)

		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, created.ID)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestPublicationRepository_Links(t *testing.T) {
	dbManager := setupTestDB(t)
	pubRepo := NewPublicationRepository(dbManager)
	memberRepo := NewLabMemberRepository(dbManager)

	t.Run("link and unlink author", func(t *testing.T) {
		// Create a publication
		pub := &models.Publication{
			Title:       "Multi-Author Paper",
			AuthorsText: "Multiple Authors",
			Year:        2024,
		}
		createdPub, err := pubRepo.Create(ctx, pub)
		require.NoError(t, err)

		// Create a lab member
		member := &models.LabMember{
			Name:         "Co Author",
			Role:         models.LabMemberRolePhD,
			DisplayOrder: 1,
		}
		createdMember, err := memberRepo.Create(ctx, member)
		require.NoError(t, err)

		// Link author
		err = pubRepo.LinkAuthor(ctx, createdPub.ID, createdMember.ID)
		require.NoError(t, err)

		// Get authors
		authors, err := pubRepo.GetAuthors(ctx, createdPub.ID)
		require.NoError(t, err)
		assert.Len(t, authors, 1)
		assert.Equal(t, createdMember.ID, authors[0].ID)

		// Unlink author
		err = pubRepo.UnlinkAuthor(ctx, createdPub.ID, createdMember.ID)
		require.NoError(t, err)

		// Verify unlink
		authors, err = pubRepo.GetAuthors(ctx, createdPub.ID)
		require.NoError(t, err)
		assert.Empty(t, authors)
	})

	t.Run("get by member", func(t *testing.T) {
		// Create member
		member := &models.LabMember{
			Name:         "Researcher",
			Role:         models.LabMemberRolePostdoc,
			DisplayOrder: 2,
		}
		createdMember, err := memberRepo.Create(ctx, member)
		require.NoError(t, err)

		// Create publication
		pub := &models.Publication{
			Title:       "Linked Paper",
			AuthorsText: "Authors",
			Year:        2023,
		}
		createdPub, err := pubRepo.Create(ctx, pub)
		require.NoError(t, err)

		// Link them
		err = pubRepo.LinkAuthor(ctx, createdPub.ID, createdMember.ID)
		require.NoError(t, err)

		// Get publications by member
		pubs, err := pubRepo.GetByMember(ctx, createdMember.ID)
		require.NoError(t, err)
		assert.NotEmpty(t, pubs)
	})

	t.Run("get publication with authors", func(t *testing.T) {
		// Create publication
		pub := &models.Publication{
			Title:       "Full Paper",
			AuthorsText: "Authors",
			Year:        2022,
		}
		createdPub, err := pubRepo.Create(ctx, pub)
		require.NoError(t, err)

		// Create member
		member := &models.LabMember{
			Name:         "Full Author",
			Role:         models.LabMemberRolePI,
			DisplayOrder: 3,
		}
		createdMember, err := memberRepo.Create(ctx, member)
		require.NoError(t, err)

		// Link them
		err = pubRepo.LinkAuthor(ctx, createdPub.ID, createdMember.ID)
		require.NoError(t, err)

		// Get with authors
		pubWithAuthors, err := pubRepo.GetWithAuthors(ctx, createdPub.ID)
		require.NoError(t, err)
		assert.Equal(t, createdPub.ID, pubWithAuthors.ID)
		assert.Len(t, pubWithAuthors.Authors, 1)
	})
}
