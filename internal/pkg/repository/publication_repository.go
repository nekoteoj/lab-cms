package repository

import (
	"context"
	"database/sql"

	"github.com/nekoteoj/lab-cms/internal/pkg/db"
	"github.com/nekoteoj/lab-cms/internal/pkg/models"
)

// Ensure PublicationRepository implements Repository[Publication] interface
var _ Repository[models.Publication] = (*PublicationRepository)(nil)

// PublicationRepository provides data access for publications.
type PublicationRepository struct {
	*BaseRepository
}

// NewPublicationRepository creates a new publication repository.
func NewPublicationRepository(dbManager *db.DBManager) *PublicationRepository {
	return &PublicationRepository{
		BaseRepository: NewBaseRepository(dbManager, "publications"),
	}
}

// GetByID retrieves a publication by ID.
func (r *PublicationRepository) GetByID(ctx context.Context, id int) (*models.Publication, error) {
	query := `
		SELECT id, title, authors_text, venue, year, url, created_at, updated_at
		FROM publications
		WHERE id = $1
	`

	row := r.GetExecer(ctx).QueryRowContext(ctx, query, id)

	var pub models.Publication
	err := row.Scan(
		&pub.ID,
		&pub.Title,
		&pub.AuthorsText,
		&pub.Venue,
		&pub.Year,
		&pub.URL,
		&pub.CreatedAt,
		&pub.UpdatedAt,
	)

	if err != nil {
		return nil, WrapError(err, "get publication by id")
	}

	return &pub, nil
}

// GetAll retrieves all publications ordered by year (newest first).
func (r *PublicationRepository) GetAll(ctx context.Context) ([]models.Publication, error) {
	query := `
		SELECT id, title, authors_text, venue, year, url, created_at, updated_at
		FROM publications
		ORDER BY year DESC, created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, WrapError(err, "get all publications")
	}
	defer rows.Close()

	var pubs []models.Publication
	for rows.Next() {
		var pub models.Publication
		err := rows.Scan(
			&pub.ID,
			&pub.Title,
			&pub.AuthorsText,
			&pub.Venue,
			&pub.Year,
			&pub.URL,
			&pub.CreatedAt,
			&pub.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan publication")
		}
		pubs = append(pubs, pub)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate publications")
	}

	return pubs, nil
}

// GetByYear retrieves publications for a specific year.
func (r *PublicationRepository) GetByYear(ctx context.Context, year int) ([]models.Publication, error) {
	query := `
		SELECT id, title, authors_text, venue, year, url, created_at, updated_at
		FROM publications
		WHERE year = $1
		ORDER BY created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query, year)
	if err != nil {
		return nil, WrapError(err, "get publications by year")
	}
	defer rows.Close()

	var pubs []models.Publication
	for rows.Next() {
		var pub models.Publication
		err := rows.Scan(
			&pub.ID,
			&pub.Title,
			&pub.AuthorsText,
			&pub.Venue,
			&pub.Year,
			&pub.URL,
			&pub.CreatedAt,
			&pub.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan publication")
		}
		pubs = append(pubs, pub)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate publications by year")
	}

	return pubs, nil
}

// GetByMember retrieves publications associated with a lab member.
func (r *PublicationRepository) GetByMember(ctx context.Context, memberID int) ([]models.Publication, error) {
	query := `
		SELECT p.id, p.title, p.authors_text, p.venue, p.year, p.url, p.created_at, p.updated_at
		FROM publications p
		INNER JOIN publication_authors pa ON p.id = pa.publication_id
		WHERE pa.member_id = $1
		ORDER BY p.year DESC, p.created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query, memberID)
	if err != nil {
		return nil, WrapError(err, "get publications by member")
	}
	defer rows.Close()

	var pubs []models.Publication
	for rows.Next() {
		var pub models.Publication
		err := rows.Scan(
			&pub.ID,
			&pub.Title,
			&pub.AuthorsText,
			&pub.Venue,
			&pub.Year,
			&pub.URL,
			&pub.CreatedAt,
			&pub.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan publication")
		}
		pubs = append(pubs, pub)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate publications by member")
	}

	return pubs, nil
}

// Create inserts a new publication.
func (r *PublicationRepository) Create(ctx context.Context, pub *models.Publication) (*models.Publication, error) {
	query := `
		INSERT INTO publications (title, authors_text, venue, year, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, datetime('now'), datetime('now'))
		RETURNING id, created_at, updated_at
	`

	row := r.GetExecer(ctx).QueryRowContext(
		ctx,
		query,
		pub.Title,
		pub.AuthorsText,
		pub.Venue,
		pub.Year,
		pub.URL,
	)

	err := row.Scan(&pub.ID, &pub.CreatedAt, &pub.UpdatedAt)
	if err != nil {
		return nil, WrapError(err, "create publication")
	}

	return pub, nil
}

// Update modifies an existing publication.
func (r *PublicationRepository) Update(ctx context.Context, pub *models.Publication) (*models.Publication, error) {
	query := `
		UPDATE publications
		SET title = $1, authors_text = $2, venue = $3, year = $4, url = $5,
		    updated_at = datetime('now')
		WHERE id = $6
		RETURNING updated_at
	`

	row := r.GetExecer(ctx).QueryRowContext(
		ctx,
		query,
		pub.Title,
		pub.AuthorsText,
		pub.Venue,
		pub.Year,
		pub.URL,
		pub.ID,
	)

	err := row.Scan(&pub.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, WrapError(err, "update publication")
	}

	return pub, nil
}

// Delete removes a publication.
func (r *PublicationRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM publications WHERE id = $1`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return WrapError(err, "delete publication")
	}

	return CheckRowsAffected(result, 1)
}

// LinkAuthor associates a lab member with a publication.
func (r *PublicationRepository) LinkAuthor(ctx context.Context, publicationID, memberID int) error {
	query := `
		INSERT INTO publication_authors (publication_id, member_id)
		VALUES ($1, $2)
		ON CONFLICT (publication_id, member_id) DO NOTHING
	`

	_, err := r.GetExecer(ctx).ExecContext(ctx, query, publicationID, memberID)
	if err != nil {
		return WrapError(err, "link author to publication")
	}

	return nil
}

// UnlinkAuthor removes the association between a lab member and a publication.
func (r *PublicationRepository) UnlinkAuthor(ctx context.Context, publicationID, memberID int) error {
	query := `DELETE FROM publication_authors WHERE publication_id = $1 AND member_id = $2`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, publicationID, memberID)
	if err != nil {
		return WrapError(err, "unlink author from publication")
	}

	return CheckRowsAffected(result, 1)
}

// GetAuthors retrieves all authors for a publication.
func (r *PublicationRepository) GetAuthors(ctx context.Context, publicationID int) ([]models.LabMember, error) {
	query := `
		SELECT m.id, m.name, m.role, m.email, m.bio, m.photo_url,
		       m.personal_page_content, m.research_interests, m.is_alumni,
		       m.display_order, m.created_at, m.updated_at
		FROM lab_members m
		INNER JOIN publication_authors pa ON m.id = pa.member_id
		WHERE pa.publication_id = $1
		ORDER BY m.display_order ASC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query, publicationID)
	if err != nil {
		return nil, WrapError(err, "get publication authors")
	}
	defer rows.Close()

	var members []models.LabMember
	for rows.Next() {
		var m models.LabMember
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Role,
			&m.Email,
			&m.Bio,
			&m.PhotoURL,
			&m.PersonalPageContent,
			&m.ResearchInterests,
			&m.IsAlumni,
			&m.DisplayOrder,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan author")
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate authors")
	}

	return members, nil
}

// GetWithAuthors retrieves a publication with its authors.
func (r *PublicationRepository) GetWithAuthors(ctx context.Context, id int) (*models.PublicationWithAuthors, error) {
	pub, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	authors, err := r.GetAuthors(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.PublicationWithAuthors{
		Publication: *pub,
		Authors:     authors,
	}, nil
}
