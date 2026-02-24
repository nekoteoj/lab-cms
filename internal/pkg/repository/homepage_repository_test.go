package repository

import (
	"testing"

	"github.com/nekoteoj/lab-cms/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHomepageRepository_CRUD(t *testing.T) {
	dbManager := setupTestDB(t)
	repo := NewHomepageRepository(dbManager)

	t.Run("create section", func(t *testing.T) {
		section := &models.HomepageSection{
			SectionKey:   "test_section",
			Title:        "Test Section",
			Content:      "This is test content for the section.",
			DisplayOrder: 1,
		}

		created, err := repo.Create(ctx, section)
		require.NoError(t, err)
		assert.Greater(t, created.ID, 0)
		assert.Equal(t, "test_section", created.SectionKey)
		assert.Equal(t, "Test Section", created.Title)
	})

	t.Run("get section by id", func(t *testing.T) {
		section := &models.HomepageSection{
			SectionKey:   "get_test",
			Title:        "Get Test",
			Content:      "Content for get test",
			DisplayOrder: 2,
		}

		created, err := repo.Create(ctx, section)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, "get_test", retrieved.SectionKey)
	})

	t.Run("get section by key", func(t *testing.T) {
		section := &models.HomepageSection{
			SectionKey:   "key_test",
			Title:        "Key Test",
			Content:      "Content for key test",
			DisplayOrder: 3,
		}

		created, err := repo.Create(ctx, section)
		require.NoError(t, err)

		retrieved, err := repo.GetByKey(ctx, "key_test")
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, "key_test", retrieved.SectionKey)
	})

	t.Run("get all sections", func(t *testing.T) {
		// Create multiple sections
		for i := 0; i < 3; i++ {
			section := &models.HomepageSection{
				SectionKey:   "section_" + string(rune('a'+i)),
				Title:        "Section " + string(rune('A'+i)),
				Content:      "Content",
				DisplayOrder: i + 10,
			}
			_, err := repo.Create(ctx, section)
			require.NoError(t, err)
		}

		sections, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, sections)
	})

	t.Run("update section", func(t *testing.T) {
		section := &models.HomepageSection{
			SectionKey:   "update_test",
			Title:        "Original Title",
			Content:      "Original content",
			DisplayOrder: 5,
		}

		created, err := repo.Create(ctx, section)
		require.NoError(t, err)

		created.Title = "Updated Title"
		created.Content = "Updated content"

		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, "Updated content", updated.Content)
	})

	t.Run("update content", func(t *testing.T) {
		section := &models.HomepageSection{
			SectionKey:   "content_test",
			Title:        "Content Test",
			Content:      "Original content",
			DisplayOrder: 6,
		}

		created, err := repo.Create(ctx, section)
		require.NoError(t, err)

		err = repo.UpdateContent(ctx, created.ID, "New Title", "New Content")
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, "New Title", retrieved.Title)
		assert.Equal(t, "New Content", retrieved.Content)
	})

	t.Run("update content by key", func(t *testing.T) {
		section := &models.HomepageSection{
			SectionKey:   "key_update_test",
			Title:        "Original",
			Content:      "Original",
			DisplayOrder: 7,
		}

		created, err := repo.Create(ctx, section)
		require.NoError(t, err)

		err = repo.UpdateContentByKey(ctx, "key_update_test", "Updated By Key", "Content Updated")
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated By Key", retrieved.Title)
		assert.Equal(t, "Content Updated", retrieved.Content)
	})

	t.Run("delete section", func(t *testing.T) {
		section := &models.HomepageSection{
			SectionKey:   "delete_test",
			Title:        "To Delete",
			Content:      "Will be deleted",
			DisplayOrder: 8,
		}

		created, err := repo.Create(ctx, section)
		require.NoError(t, err)

		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, created.ID)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("duplicate key error", func(t *testing.T) {
		section1 := &models.HomepageSection{
			SectionKey:   "duplicate_key",
			Title:        "First",
			Content:      "Content",
			DisplayOrder: 9,
		}

		_, err := repo.Create(ctx, section1)
		require.NoError(t, err)

		section2 := &models.HomepageSection{
			SectionKey:   "duplicate_key",
			Title:        "Second",
			Content:      "Content",
			DisplayOrder: 10,
		}

		_, err = repo.Create(ctx, section2)
		assert.Equal(t, ErrDuplicate, err)
	})
}
