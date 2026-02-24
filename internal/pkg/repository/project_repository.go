package repository

import (
	"context"
	"database/sql"

	"github.com/nekoteoj/lab-cms/internal/pkg/db"
	"github.com/nekoteoj/lab-cms/internal/pkg/models"
)

// Ensure ProjectRepository implements Repository[Project] interface
var _ Repository[models.Project] = (*ProjectRepository)(nil)

// ProjectRepository provides data access for projects.
type ProjectRepository struct {
	*BaseRepository
}

// NewProjectRepository creates a new project repository.
func NewProjectRepository(dbManager *db.DBManager) *ProjectRepository {
	return &ProjectRepository{
		BaseRepository: NewBaseRepository(dbManager, "projects"),
	}
}

// GetByID retrieves a project by ID.
func (r *ProjectRepository) GetByID(ctx context.Context, id int) (*models.Project, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	row := r.GetExecer(ctx).QueryRowContext(ctx, query, id)

	var proj models.Project
	err := row.Scan(
		&proj.ID,
		&proj.Title,
		&proj.Description,
		&proj.Status,
		&proj.CreatedAt,
		&proj.UpdatedAt,
	)

	if err != nil {
		return nil, WrapError(err, "get project by id")
	}

	return &proj, nil
}

// GetAll retrieves all projects ordered by status and creation date.
func (r *ProjectRepository) GetAll(ctx context.Context) ([]models.Project, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM projects
		ORDER BY 
			CASE status WHEN 'active' THEN 0 ELSE 1 END,
			created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, WrapError(err, "get all projects")
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var proj models.Project
		err := rows.Scan(
			&proj.ID,
			&proj.Title,
			&proj.Description,
			&proj.Status,
			&proj.CreatedAt,
			&proj.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan project")
		}
		projects = append(projects, proj)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate projects")
	}

	return projects, nil
}

// GetByStatus retrieves projects filtered by status.
func (r *ProjectRepository) GetByStatus(ctx context.Context, status models.ProjectStatus) ([]models.Project, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM projects
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query, status)
	if err != nil {
		return nil, WrapError(err, "get projects by status")
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var proj models.Project
		err := rows.Scan(
			&proj.ID,
			&proj.Title,
			&proj.Description,
			&proj.Status,
			&proj.CreatedAt,
			&proj.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan project")
		}
		projects = append(projects, proj)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate projects by status")
	}

	return projects, nil
}

// Create inserts a new project.
func (r *ProjectRepository) Create(ctx context.Context, proj *models.Project) (*models.Project, error) {
	query := `
		INSERT INTO projects (title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, datetime('now'), datetime('now'))
		RETURNING id, created_at, updated_at
	`

	row := r.GetExecer(ctx).QueryRowContext(
		ctx,
		query,
		proj.Title,
		proj.Description,
		proj.Status,
	)

	err := row.Scan(&proj.ID, &proj.CreatedAt, &proj.UpdatedAt)
	if err != nil {
		return nil, WrapError(err, "create project")
	}

	return proj, nil
}

// Update modifies an existing project.
func (r *ProjectRepository) Update(ctx context.Context, proj *models.Project) (*models.Project, error) {
	query := `
		UPDATE projects
		SET title = $1, description = $2, status = $3, updated_at = datetime('now')
		WHERE id = $4
		RETURNING updated_at
	`

	row := r.GetExecer(ctx).QueryRowContext(
		ctx,
		query,
		proj.Title,
		proj.Description,
		proj.Status,
		proj.ID,
	)

	err := row.Scan(&proj.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, WrapError(err, "update project")
	}

	return proj, nil
}

// Delete removes a project.
func (r *ProjectRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM projects WHERE id = $1`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return WrapError(err, "delete project")
	}

	return CheckRowsAffected(result, 1)
}

// LinkMember associates a lab member with a project.
func (r *ProjectRepository) LinkMember(ctx context.Context, projectID, memberID int) error {
	query := `
		INSERT INTO project_members (project_id, member_id)
		VALUES ($1, $2)
		ON CONFLICT (project_id, member_id) DO NOTHING
	`

	_, err := r.GetExecer(ctx).ExecContext(ctx, query, projectID, memberID)
	if err != nil {
		return WrapError(err, "link member to project")
	}

	return nil
}

// UnlinkMember removes the association between a lab member and a project.
func (r *ProjectRepository) UnlinkMember(ctx context.Context, projectID, memberID int) error {
	query := `DELETE FROM project_members WHERE project_id = $1 AND member_id = $2`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, projectID, memberID)
	if err != nil {
		return WrapError(err, "unlink member from project")
	}

	return CheckRowsAffected(result, 1)
}

// LinkPublication associates a publication with a project.
func (r *ProjectRepository) LinkPublication(ctx context.Context, projectID, publicationID int) error {
	query := `
		INSERT INTO project_publications (project_id, publication_id)
		VALUES ($1, $2)
		ON CONFLICT (project_id, publication_id) DO NOTHING
	`

	_, err := r.GetExecer(ctx).ExecContext(ctx, query, projectID, publicationID)
	if err != nil {
		return WrapError(err, "link publication to project")
	}

	return nil
}

// UnlinkPublication removes the association between a publication and a project.
func (r *ProjectRepository) UnlinkPublication(ctx context.Context, projectID, publicationID int) error {
	query := `DELETE FROM project_publications WHERE project_id = $1 AND publication_id = $2`

	result, err := r.GetExecer(ctx).ExecContext(ctx, query, projectID, publicationID)
	if err != nil {
		return WrapError(err, "unlink publication from project")
	}

	return CheckRowsAffected(result, 1)
}

// GetMembers retrieves all members associated with a project.
func (r *ProjectRepository) GetMembers(ctx context.Context, projectID int) ([]models.LabMember, error) {
	query := `
		SELECT m.id, m.name, m.role, m.email, m.bio, m.photo_url,
		       m.personal_page_content, m.research_interests, m.is_alumni,
		       m.display_order, m.created_at, m.updated_at
		FROM lab_members m
		INNER JOIN project_members pm ON m.id = pm.member_id
		WHERE pm.project_id = $1
		ORDER BY m.display_order ASC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, WrapError(err, "get project members")
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
			return nil, WrapError(err, "scan project member")
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate project members")
	}

	return members, nil
}

// GetPublications retrieves all publications associated with a project.
func (r *ProjectRepository) GetPublications(ctx context.Context, projectID int) ([]models.Publication, error) {
	query := `
		SELECT p.id, p.title, p.authors_text, p.venue, p.year, p.url, p.created_at, p.updated_at
		FROM publications p
		INNER JOIN project_publications pp ON p.id = pp.publication_id
		WHERE pp.project_id = $1
		ORDER BY p.year DESC
	`

	rows, err := r.GetExecer(ctx).QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, WrapError(err, "get project publications")
	}
	defer rows.Close()

	var pubs []models.Publication
	for rows.Next() {
		var p models.Publication
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.AuthorsText,
			&p.Venue,
			&p.Year,
			&p.URL,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, WrapError(err, "scan project publication")
		}
		pubs = append(pubs, p)
	}

	if err := rows.Err(); err != nil {
		return nil, WrapError(err, "iterate project publications")
	}

	return pubs, nil
}

// GetWithRelations retrieves a project with its members and publications.
func (r *ProjectRepository) GetWithRelations(ctx context.Context, id int) (*models.ProjectWithRelations, error) {
	proj, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	members, err := r.GetMembers(ctx, id)
	if err != nil {
		return nil, err
	}

	publications, err := r.GetPublications(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.ProjectWithRelations{
		Project:      *proj,
		Members:      members,
		Publications: publications,
	}, nil
}
