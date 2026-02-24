package repository

import (
	"context"
	"database/sql"

	"github.com/nekoteoj/lab-cms/internal/pkg/db"
	"github.com/nekoteoj/lab-cms/internal/pkg/models"
)

// Ensure NewsRepository implements Repository[News] interface
var _ Repository[models.News] = (*NewsRepository)(nil)

// NewsRepository provides data access for news items.
type NewsRepository struct {
	*BaseRepository
}

// NewNewsRepository creates a new news repository.
func NewNewsRepository(dbManager *db.DBManager) *NewsRepository {
	return &NewsRepository{
		BaseRepository: NewBaseRepository(dbManager, "news"),
	}
}

// GetByID retrieves a news item by ID.
func (r *NewsRepository) GetByID(ctx context.Context, id int) (*models.News, error) {
	query := `
		SELECT id, title, content, published_at, is_published, created_at, updated_at
		FROM news
		WHERE id = $1
	`

	row := r.GetExecer(ctx).QueryRowContext(ctx, query, id)

	var news models.News
	err := row.Scan(
		&news.ID,
		&news.Title,
		&news.Content,
		&news.PublishedAt,
		&news.IsPublished,
		&news.CreatedAt,
		&news.UpdatedAt,
	)

	if err != nil {
		return nil, WrapError(err, "get news by id")
	}

	return &news, nil
}

// GetAll retrieves all news items ordered by creation date.
func (r *NewsRepository) GetAll(ctx context.Context) ([]models.News, error) {
	query := `
		SELECT id, title, content, published_at, is_published, created_at, updated_at
		FROM news
		ORDER BY created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, WrapError(err, "get all news")
	}
	defer rows.Close()

	var news []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&n.PublishedAt,
			&n.IsPublished,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan news")
		}
		news = append(news, n)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate news")
	}

	return news, nil
}

// GetPublished retrieves all published news items that should be visible to the public.
func (r *NewsRepository) GetPublished(ctx context.Context, limit int) ([]models.News, error) {
	query := `
		SELECT id, title, content, published_at, is_published, created_at, updated_at
		FROM news
		WHERE is_published = true
		  AND (published_at IS NULL OR published_at <= datetime('now'))
		ORDER BY 
			CASE WHEN published_at IS NOT NULL THEN published_at ELSE created_at END DESC
		LIMIT $1
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query, limit)
	if err != nil {
		return nil, WrapError(err, "get published news")
	}
	defer rows.Close()

	var news []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&n.PublishedAt,
			&n.IsPublished,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan news")
		}
		news = append(news, n)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate published news")
	}

	return news, nil
}

// GetDrafts retrieves all unpublished news items.
func (r *NewsRepository) GetDrafts(ctx context.Context) ([]models.News, error) {
	query := `
		SELECT id, title, content, published_at, is_published, created_at, updated_at
		FROM news
		WHERE is_published = false
		ORDER BY created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, WrapError(err, "get draft news")
	}
	defer rows.Close()

	var news []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&n.PublishedAt,
			&n.IsPublished,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan news")
		}
		news = append(news, n)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate draft news")
	}

	return news, nil
}

// Create inserts a new news item.
func (r *NewsRepository) Create(ctx context.Context, news *models.News) (*models.News, error) {
	var query string
	var row *sql.Row

	if news.PublishedAt.Valid {
		// News with specific publish date
		query = `
			INSERT INTO news (title, content, published_at, is_published, created_at, updated_at)
			VALUES ($1, $2, $3, $4, datetime('now'), datetime('now'))
			RETURNING id, created_at, updated_at
		`
		row = r.GetExecer(ctx).QueryRowContext(
			ctx,
			query,
			news.Title,
			news.Content,
			news.PublishedAt,
			news.IsPublished,
		)
	} else {
		// News without specific publish date
		query = `
			INSERT INTO news (title, content, published_at, is_published, created_at, updated_at)
			VALUES ($1, $2, NULL, $3, datetime('now'), datetime('now'))
			RETURNING id, created_at, updated_at
		`
		row = r.GetExecer(ctx).QueryRowContext(
			ctx,
			query,
			news.Title,
			news.Content,
			news.IsPublished,
		)
	}

	err := row.Scan(&news.ID, &news.CreatedAt, &news.UpdatedAt)
	if err != nil {
		return nil, WrapError(err, "create news")
	}

	return news, nil
}

// Update modifies an existing news item.
func (r *NewsRepository) Update(ctx context.Context, news *models.News) (*models.News, error) {
	var query string
	var row *sql.Row

	if news.PublishedAt.Valid {
		query = `
			UPDATE news
			SET title = $1, content = $2, published_at = $3, is_published = $4,
			    updated_at = datetime('now')
			WHERE id = $5
			RETURNING updated_at
		`
		row = r.GetExecer(ctx).QueryRowContext(
			ctx,
			query,
			news.Title,
			news.Content,
			news.PublishedAt,
			news.IsPublished,
			news.ID,
		)
	} else {
		query = `
			UPDATE news
			SET title = $1, content = $2, published_at = NULL, is_published = $3,
			    updated_at = datetime('now')
			WHERE id = $4
			RETURNING updated_at
		`
		row = r.GetExecer(ctx).QueryRowContext(
			ctx,
			query,
			news.Title,
			news.Content,
			news.IsPublished,
			news.ID,
		)
	}

	err := row.Scan(&news.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, WrapError(err, "update news")
	}

	return news, nil
}

// Delete removes a news item.
func (r *NewsRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM news WHERE id = $1`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return WrapError(err, "delete news")
	}

	return CheckRowsAffected(result, 1)
}

// Publish marks a news item as published.
func (r *NewsRepository) Publish(ctx context.Context, id int) error {
	query := `
		UPDATE news
		SET is_published = true, published_at = datetime('now'), updated_at = datetime('now')
		WHERE id = $1
	`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return WrapError(err, "publish news")
	}

	return CheckRowsAffected(result, 1)
}

// Unpublish marks a news item as unpublished.
func (r *NewsRepository) Unpublish(ctx context.Context, id int) error {
	query := `
		UPDATE news
		SET is_published = false, updated_at = datetime('now')
		WHERE id = $1
	`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return WrapError(err, "unpublish news")
	}

	return CheckRowsAffected(result, 1)
}
