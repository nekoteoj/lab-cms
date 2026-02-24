package repository

import (
	"context"
	"database/sql"

	"github.com/nekoteoj/lab-cms/internal/pkg/db"
	"github.com/nekoteoj/lab-cms/internal/pkg/models"
)

// Ensure LabMemberRepository implements Repository[LabMember] interface
var _ Repository[models.LabMember] = (*LabMemberRepository)(nil)

// LabMemberRepository provides data access for lab members.
type LabMemberRepository struct {
	*BaseRepository
}

// NewLabMemberRepository creates a new lab member repository.
func NewLabMemberRepository(dbManager *db.DBManager) *LabMemberRepository {
	return &LabMemberRepository{
		BaseRepository: NewBaseRepository(dbManager, "lab_members"),
	}
}

// GetByID retrieves a lab member by ID.
func (r *LabMemberRepository) GetByID(ctx context.Context, id int) (*models.LabMember, error) {
	query := `
		SELECT id, name, role, email, bio, photo_url, personal_page_content,
		       research_interests, is_alumni, display_order, created_at, updated_at
		FROM lab_members
		WHERE id = $1
	`

	row := r.GetExecer(ctx).QueryRowContext(ctx, query, id)

	var member models.LabMember
	err := row.Scan(
		&member.ID,
		&member.Name,
		&member.Role,
		&member.Email,
		&member.Bio,
		&member.PhotoURL,
		&member.PersonalPageContent,
		&member.ResearchInterests,
		&member.IsAlumni,
		&member.DisplayOrder,
		&member.CreatedAt,
		&member.UpdatedAt,
	)

	if err != nil {
		return nil, WrapError(err, "get lab member by id")
	}

	return &member, nil
}

// GetAll retrieves all lab members ordered by display_order.
func (r *LabMemberRepository) GetAll(ctx context.Context) ([]models.LabMember, error) {
	query := `
		SELECT id, name, role, email, bio, photo_url, personal_page_content,
		       research_interests, is_alumni, display_order, created_at, updated_at
		FROM lab_members
		ORDER BY is_alumni ASC, display_order ASC, created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, WrapError(err, "get all lab members")
	}
	defer rows.Close()

	var members []models.LabMember
	for rows.Next() {
		var member models.LabMember
		err := rows.Scan(
			&member.ID,
			&member.Name,
			&member.Role,
			&member.Email,
			&member.Bio,
			&member.PhotoURL,
			&member.PersonalPageContent,
			&member.ResearchInterests,
			&member.IsAlumni,
			&member.DisplayOrder,
			&member.CreatedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan lab member")
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate lab members")
	}

	return members, nil
}

// GetByRole retrieves lab members filtered by role.
func (r *LabMemberRepository) GetByRole(ctx context.Context, role models.LabMemberRole) ([]models.LabMember, error) {
	query := `
		SELECT id, name, role, email, bio, photo_url, personal_page_content,
		       research_interests, is_alumni, display_order, created_at, updated_at
		FROM lab_members
		WHERE role = $1 AND is_alumni = false
		ORDER BY display_order ASC, created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query, role)
	if err != nil {
		return nil, WrapError(err, "get lab members by role")
	}
	defer rows.Close()

	var members []models.LabMember
	for rows.Next() {
		var member models.LabMember
		err := rows.Scan(
			&member.ID,
			&member.Name,
			&member.Role,
			&member.Email,
			&member.Bio,
			&member.PhotoURL,
			&member.PersonalPageContent,
			&member.ResearchInterests,
			&member.IsAlumni,
			&member.DisplayOrder,
			&member.CreatedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan lab member")
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate lab members by role")
	}

	return members, nil
}

// GetAlumni retrieves all alumni members.
func (r *LabMemberRepository) GetAlumni(ctx context.Context) ([]models.LabMember, error) {
	query := `
		SELECT id, name, role, email, bio, photo_url, personal_page_content,
		       research_interests, is_alumni, display_order, created_at, updated_at
		FROM lab_members
		WHERE is_alumni = true
		ORDER BY display_order ASC, created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, WrapError(err, "get alumni")
	}
	defer rows.Close()

	var members []models.LabMember
	for rows.Next() {
		var member models.LabMember
		err := rows.Scan(
			&member.ID,
			&member.Name,
			&member.Role,
			&member.Email,
			&member.Bio,
			&member.PhotoURL,
			&member.PersonalPageContent,
			&member.ResearchInterests,
			&member.IsAlumni,
			&member.DisplayOrder,
			&member.CreatedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan alumni member")
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate alumni")
	}

	return members, nil
}

// Create inserts a new lab member.
func (r *LabMemberRepository) Create(ctx context.Context, member *models.LabMember) (*models.LabMember, error) {
	query := `
		INSERT INTO lab_members (
			name, role, email, bio, photo_url, personal_page_content,
			research_interests, is_alumni, display_order, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			datetime('now'), datetime('now')
		)
		RETURNING id, created_at, updated_at
	`

	row := r.GetExecer(ctx).QueryRowContext(
		ctx,
		query,
		member.Name,
		member.Role,
		member.Email,
		member.Bio,
		member.PhotoURL,
		member.PersonalPageContent,
		member.ResearchInterests,
		member.IsAlumni,
		member.DisplayOrder,
	)

	err := row.Scan(&member.ID, &member.CreatedAt, &member.UpdatedAt)
	if err != nil {
		return nil, WrapError(err, "create lab member")
	}

	return member, nil
}

// Update modifies an existing lab member.
func (r *LabMemberRepository) Update(ctx context.Context, member *models.LabMember) (*models.LabMember, error) {
	query := `
		UPDATE lab_members
		SET name = $1, role = $2, email = $3, bio = $4, photo_url = $5,
		    personal_page_content = $6, research_interests = $7, is_alumni = $8,
		    display_order = $9, updated_at = datetime('now')
		WHERE id = $10
		RETURNING updated_at
	`

	row := r.GetExecer(ctx).QueryRowContext(
		ctx,
		query,
		member.Name,
		member.Role,
		member.Email,
		member.Bio,
		member.PhotoURL,
		member.PersonalPageContent,
		member.ResearchInterests,
		member.IsAlumni,
		member.DisplayOrder,
		member.ID,
	)

	err := row.Scan(&member.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, WrapError(err, "update lab member")
	}

	return member, nil
}

// Delete removes a lab member.
func (r *LabMemberRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM lab_members WHERE id = $1`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return WrapError(err, "delete lab member")
	}

	return CheckRowsAffected(result, 1)
}

// MarkAsAlumni updates a member's alumni status.
func (r *LabMemberRepository) MarkAsAlumni(ctx context.Context, id int, isAlumni bool) error {
	query := `
		UPDATE lab_members
		SET is_alumni = $1, updated_at = datetime('now')
		WHERE id = $2
	`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, isAlumni, id)
	if err != nil {
		return WrapError(err, "mark member as alumni")
	}

	return CheckRowsAffected(result, 1)
}

// UpdatePhotoURL updates a member's photo URL.
func (r *LabMemberRepository) UpdatePhotoURL(ctx context.Context, id int, photoURL string) error {
	query := `
		UPDATE lab_members
		SET photo_url = $1, updated_at = datetime('now')
		WHERE id = $2
	`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, photoURL, id)
	if err != nil {
		return WrapError(err, "update member photo")
	}

	return CheckRowsAffected(result, 1)
}
