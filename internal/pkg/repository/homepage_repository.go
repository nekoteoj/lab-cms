package repository

import (
	"context"
	"database/sql"

	"github.com/nekoteoj/lab-cms/internal/pkg/db"
	"github.com/nekoteoj/lab-cms/internal/pkg/models"
)

// Ensure HomepageRepository implements Repository[HomepageSection] interface
var _ Repository[models.HomepageSection] = (*HomepageRepository)(nil)

// HomepageRepository provides data access for homepage sections.
type HomepageRepository struct {
	*BaseRepository
}

// NewHomepageRepository creates a new homepage repository.
func NewHomepageRepository(dbManager *db.DBManager) *HomepageRepository {
	return &HomepageRepository{
		BaseRepository: NewBaseRepository(dbManager, "homepage_sections"),
	}
}

// GetByID retrieves a homepage section by ID.
func (r *HomepageRepository) GetByID(ctx context.Context, id int) (*models.HomepageSection, error) {
	query := `
		SELECT id, section_key, title, content, display_order, updated_at
		FROM homepage_sections
		WHERE id = $1
	`

	row := r.GetExecer(ctx).QueryRowContext(ctx, query, id)

	var section models.HomepageSection
	err := row.Scan(
		&section.ID,
		&section.SectionKey,
		&section.Title,
		&section.Content,
		&section.DisplayOrder,
		&section.UpdatedAt,
	)

	if err != nil {
		return nil, WrapError(err, "get homepage section by id")
	}

	return &section, nil
}

// GetByKey retrieves a homepage section by its unique section key.
func (r *HomepageRepository) GetByKey(ctx context.Context, key string) (*models.HomepageSection, error) {
	query := `
		SELECT id, section_key, title, content, display_order, updated_at
		FROM homepage_sections
		WHERE section_key = $1
	`

	row := r.GetExecer(ctx).QueryRowContext(ctx, query, key)

	var section models.HomepageSection
	err := row.Scan(
		&section.ID,
		&section.SectionKey,
		&section.Title,
		&section.Content,
		&section.DisplayOrder,
		&section.UpdatedAt,
	)

	if err != nil {
		return nil, WrapError(err, "get homepage section by key")
	}

	return &section, nil
}

// GetAll retrieves all homepage sections ordered by display order.
func (r *HomepageRepository) GetAll(ctx context.Context) ([]models.HomepageSection, error) {
	query := `
		SELECT id, section_key, title, content, display_order, updated_at
		FROM homepage_sections
		ORDER BY display_order ASC, id ASC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, WrapError(err, "get all homepage sections")
	}
	defer rows.Close()

	var sections []models.HomepageSection
	for rows.Next() {
		var s models.HomepageSection
		err := rows.Scan(
			&s.ID,
			&s.SectionKey,
			&s.Title,
			&s.Content,
			&s.DisplayOrder,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan homepage section")
		}
		sections = append(sections, s)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate homepage sections")
	}

	return sections, nil
}

// Create inserts a new homepage section.
// Note: In practice, sections are typically seeded at initialization,
// but this method allows dynamic creation if needed.
func (r *HomepageRepository) Create(ctx context.Context, section *models.HomepageSection) (*models.HomepageSection, error) {
	query := `
		INSERT INTO homepage_sections (section_key, title, content, display_order, updated_at)
		VALUES ($1, $2, $3, $4, datetime('now'))
		RETURNING id, updated_at
	`

	row := r.GetExecer(ctx).QueryRowContext(
		ctx,
		query,
		section.SectionKey,
		section.Title,
		section.Content,
		section.DisplayOrder,
	)

	err := row.Scan(&section.ID, &section.UpdatedAt)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, ErrDuplicate
		}
		return nil, WrapError(err, "create homepage section")
	}

	return section, nil
}

// Update modifies an existing homepage section.
func (r *HomepageRepository) Update(ctx context.Context, section *models.HomepageSection) (*models.HomepageSection, error) {
	query := `
		UPDATE homepage_sections
		SET title = $1, content = $2, display_order = $3, updated_at = datetime('now')
		WHERE id = $4
		RETURNING updated_at
	`

	row := r.GetExecer(ctx).QueryRowContext(
		ctx,
		query,
		section.Title,
		section.Content,
		section.DisplayOrder,
		section.ID,
	)

	err := row.Scan(&section.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, WrapError(err, "update homepage section")
	}

	return section, nil
}

// Delete removes a homepage section.
// Note: Use with caution as this permanently removes the section.
func (r *HomepageRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM homepage_sections WHERE id = $1`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return WrapError(err, "delete homepage section")
	}

	return CheckRowsAffected(result, 1)
}

// UpdateContent updates just the content and title of a section.
// This is a convenience method for quick updates.
func (r *HomepageRepository) UpdateContent(ctx context.Context, id int, title, content string) error {
	query := `
		UPDATE homepage_sections
		SET title = $1, content = $2, updated_at = datetime('now')
		WHERE id = $3
	`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, title, content, id)
	if err != nil {
		return WrapError(err, "update section content")
	}

	return CheckRowsAffected(result, 1)
}

// UpdateContentByKey updates content by section key (useful for known sections like 'overview').
func (r *HomepageRepository) UpdateContentByKey(ctx context.Context, key, title, content string) error {
	query := `
		UPDATE homepage_sections
		SET title = $1, content = $2, updated_at = datetime('now')
		WHERE section_key = $3
	`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, title, content, key)
	if err != nil {
		return WrapError(err, "update section content by key")
	}

	return CheckRowsAffected(result, 1)
}

// isDuplicateKeyError checks if the error is a duplicate key violation.
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "constraint failed: UNIQUE constraint failed: homepage_sections.section_key (2067)" ||
		err.Error() == "UNIQUE constraint failed: homepage_sections.section_key"
}
