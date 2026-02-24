package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/nekoteoj/lab-cms/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewsRepository_CRUD(t *testing.T) {
	dbManager := setupTestDB(t)
	repo := NewNewsRepository(dbManager)

	t.Run("create news", func(t *testing.T) {
		news := &models.News{
			Title:       "Test News",
			Content:     "This is test content",
			IsPublished: false,
		}

		created, err := repo.Create(ctx, news)
		require.NoError(t, err)
		assert.Greater(t, created.ID, 0)
		assert.Equal(t, "Test News", created.Title)
		assert.False(t, created.IsPublished)
	})

	t.Run("get news by id", func(t *testing.T) {
		news := &models.News{
			Title:       "Another News",
			Content:     "More content",
			IsPublished: true,
			PublishedAt: sql.NullTime{Time: time.Now(), Valid: true},
		}

		created, err := repo.Create(ctx, news)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, "Another News", retrieved.Title)
	})

	t.Run("get all news", func(t *testing.T) {
		// Create multiple news items
		for i := 0; i < 3; i++ {
			news := &models.News{
				Title:       "News Item " + string(rune('A'+i)),
				Content:     "Content",
				IsPublished: i%2 == 0,
			}
			_, err := repo.Create(ctx, news)
			require.NoError(t, err)
		}

		allNews, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, allNews)
	})

	t.Run("get published news", func(t *testing.T) {
		// Create published news
		news := &models.News{
			Title:       "Published News",
			Content:     "Published content",
			IsPublished: true,
			PublishedAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
		}
		_, err := repo.Create(ctx, news)
		require.NoError(t, err)

		published, err := repo.GetPublished(ctx, 10)
		require.NoError(t, err)
		assert.NotEmpty(t, published)
	})

	t.Run("get drafts", func(t *testing.T) {
		// Create draft news
		news := &models.News{
			Title:       "Draft News",
			Content:     "Draft content",
			IsPublished: false,
		}
		_, err := repo.Create(ctx, news)
		require.NoError(t, err)

		drafts, err := repo.GetDrafts(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, drafts)
	})

	t.Run("update news", func(t *testing.T) {
		news := &models.News{
			Title:       "Original Title",
			Content:     "Original content",
			IsPublished: false,
		}

		created, err := repo.Create(ctx, news)
		require.NoError(t, err)

		created.Title = "Updated Title"
		created.Content = "Updated content"

		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, "Updated content", updated.Content)
	})

	t.Run("publish news", func(t *testing.T) {
		news := &models.News{
			Title:       "To Publish",
			Content:     "Content",
			IsPublished: false,
		}

		created, err := repo.Create(ctx, news)
		require.NoError(t, err)

		err = repo.Publish(ctx, created.ID)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.True(t, retrieved.IsPublished)
	})

	t.Run("unpublish news", func(t *testing.T) {
		news := &models.News{
			Title:       "To Unpublish",
			Content:     "Content",
			IsPublished: true,
			PublishedAt: sql.NullTime{Time: time.Now(), Valid: true},
		}

		created, err := repo.Create(ctx, news)
		require.NoError(t, err)

		err = repo.Unpublish(ctx, created.ID)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.False(t, retrieved.IsPublished)
	})

	t.Run("delete news", func(t *testing.T) {
		news := &models.News{
			Title:       "To Delete",
			Content:     "Will be deleted",
			IsPublished: false,
		}

		created, err := repo.Create(ctx, news)
		require.NoError(t, err)

		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, created.ID)
		assert.Equal(t, ErrNotFound, err)
	})
}
