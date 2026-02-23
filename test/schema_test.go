package test

import (
	"testing"

	"github.com/nekoteoj/lab-cms/test/helpers"
	"github.com/stretchr/testify/require"
)

func TestSchema_AllTablesExist(t *testing.T) {
	db := helpers.NewTestDB(t)

	tests := []struct {
		table string
	}{
		{"users"},
		{"lab_members"},
		{"publications"},
		{"projects"},
		{"news"},
		{"homepage_sections"},
		{"project_members"},
		{"publication_authors"},
		{"project_publications"},
		{"schema_migrations"},
	}

	for _, tt := range tests {
		t.Run(tt.table, func(t *testing.T) {
			exists := helpers.TableExists(t, db, tt.table)
			require.True(t, exists, "table %s should exist", tt.table)
		})
	}
}

func TestSchema_UsersTableStructure(t *testing.T) {
	db := helpers.NewTestDB(t)

	columns := helpers.GetTableColumns(t, db, "users")
	require.Contains(t, columns, "id")
	require.Contains(t, columns, "email")
	require.Contains(t, columns, "password_hash")
	require.Contains(t, columns, "role")
	require.Contains(t, columns, "created_at")
	require.Contains(t, columns, "updated_at")
}

func TestSchema_LabMembersTableStructure(t *testing.T) {
	db := helpers.NewTestDB(t)

	columns := helpers.GetTableColumns(t, db, "lab_members")
	require.Contains(t, columns, "id")
	require.Contains(t, columns, "name")
	require.Contains(t, columns, "role")
	require.Contains(t, columns, "email")
	require.Contains(t, columns, "bio")
	require.Contains(t, columns, "photo_url")
	require.Contains(t, columns, "personal_page_content")
	require.Contains(t, columns, "research_interests")
	require.Contains(t, columns, "is_alumni")
	require.Contains(t, columns, "display_order")
	require.Contains(t, columns, "created_at")
	require.Contains(t, columns, "updated_at")
}

func TestSchema_PublicationsTableStructure(t *testing.T) {
	db := helpers.NewTestDB(t)

	columns := helpers.GetTableColumns(t, db, "publications")
	require.Contains(t, columns, "id")
	require.Contains(t, columns, "title")
	require.Contains(t, columns, "authors_text")
	require.Contains(t, columns, "venue")
	require.Contains(t, columns, "year")
	require.Contains(t, columns, "url")
	require.Contains(t, columns, "created_at")
	require.Contains(t, columns, "updated_at")
}

func TestSchema_ProjectsTableStructure(t *testing.T) {
	db := helpers.NewTestDB(t)

	columns := helpers.GetTableColumns(t, db, "projects")
	require.Contains(t, columns, "id")
	require.Contains(t, columns, "title")
	require.Contains(t, columns, "description")
	require.Contains(t, columns, "status")
	require.Contains(t, columns, "created_at")
	require.Contains(t, columns, "updated_at")
}

func TestSchema_NewsTableStructure(t *testing.T) {
	db := helpers.NewTestDB(t)

	columns := helpers.GetTableColumns(t, db, "news")
	require.Contains(t, columns, "id")
	require.Contains(t, columns, "title")
	require.Contains(t, columns, "content")
	require.Contains(t, columns, "published_at")
	require.Contains(t, columns, "is_published")
	require.Contains(t, columns, "created_at")
	require.Contains(t, columns, "updated_at")
}

func TestSchema_HomepageSectionsTableStructure(t *testing.T) {
	db := helpers.NewTestDB(t)

	columns := helpers.GetTableColumns(t, db, "homepage_sections")
	require.Contains(t, columns, "id")
	require.Contains(t, columns, "section_key")
	require.Contains(t, columns, "title")
	require.Contains(t, columns, "content")
	require.Contains(t, columns, "display_order")
	require.Contains(t, columns, "updated_at")
}

func TestSchema_JunctionTableStructures(t *testing.T) {
	db := helpers.NewTestDB(t)

	tests := []struct {
		table   string
		columns []string
	}{
		{
			table:   "project_members",
			columns: []string{"project_id", "member_id"},
		},
		{
			table:   "publication_authors",
			columns: []string{"publication_id", "member_id"},
		},
		{
			table:   "project_publications",
			columns: []string{"project_id", "publication_id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.table, func(t *testing.T) {
			columns := helpers.GetTableColumns(t, db, tt.table)
			for _, col := range tt.columns {
				require.Contains(t, columns, col, "column %s should exist in %s", col, tt.table)
			}
		})
	}
}

func TestSchema_AllIndexesExist(t *testing.T) {
	db := helpers.NewTestDB(t)

	tests := []struct {
		index string
	}{
		{"idx_users_email"},
		{"idx_lab_members_alumni_order"},
		{"idx_lab_members_role"},
		{"idx_publications_year_created"},
		{"idx_projects_status"},
		{"idx_news_published_created"},
		{"idx_homepage_section_key"},
		{"idx_homepage_display_order"},
	}

	for _, tt := range tests {
		t.Run(tt.index, func(t *testing.T) {
			exists := helpers.IndexExists(t, db, tt.index)
			require.True(t, exists, "index %s should exist", tt.index)
		})
	}
}

func TestSchema_UserConstraints(t *testing.T) {
	db := helpers.NewTestDB(t)

	// Test unique email constraint
	_, err := db.Exec(
		"INSERT INTO users (email, password_hash, role) VALUES (?, ?, ?)",
		"test@example.com", "hash1", "normal",
	)
	require.NoError(t, err)

	_, err = db.Exec(
		"INSERT INTO users (email, password_hash, role) VALUES (?, ?, ?)",
		"test@example.com", "hash2", "root",
	)
	require.Error(t, err, "should error on duplicate email")
	require.Contains(t, err.Error(), "UNIQUE")

	// Test role check constraint (valid values)
	_, err = db.Exec(
		"INSERT INTO users (email, password_hash, role) VALUES (?, ?, ?)",
		"valid1@example.com", "hash", "normal",
	)
	require.NoError(t, err)

	_, err = db.Exec(
		"INSERT INTO users (email, password_hash, role) VALUES (?, ?, ?)",
		"valid2@example.com", "hash", "root",
	)
	require.NoError(t, err)

	// Test role check constraint (invalid value)
	_, err = db.Exec(
		"INSERT INTO users (email, password_hash, role) VALUES (?, ?, ?)",
		"invalid@example.com", "hash", "admin",
	)
	require.Error(t, err, "should error on invalid role")
}

func TestSchema_LabMemberConstraints(t *testing.T) {
	db := helpers.NewTestDB(t)

	// Test valid roles
	validRoles := []string{"PI", "Postdoc", "PhD", "Master", "Bachelor", "Researcher"}
	for _, role := range validRoles {
		_, err := db.Exec(
			"INSERT INTO lab_members (name, role) VALUES (?, ?)",
			"Member "+role, role,
		)
		require.NoError(t, err, "role %s should be valid", role)
	}

	// Test invalid role
	_, err := db.Exec(
		"INSERT INTO lab_members (name, role) VALUES (?, ?)",
		"Invalid Member", "InvalidRole",
	)
	require.Error(t, err, "should error on invalid member role")
}

func TestSchema_ProjectConstraints(t *testing.T) {
	db := helpers.NewTestDB(t)

	// Test valid statuses
	_, err := db.Exec(
		"INSERT INTO projects (title, description, status) VALUES (?, ?, ?)",
		"Active Project", "Description", "active",
	)
	require.NoError(t, err)

	_, err = db.Exec(
		"INSERT INTO projects (title, description, status) VALUES (?, ?, ?)",
		"Completed Project", "Description", "completed",
	)
	require.NoError(t, err)

	// Test invalid status
	_, err = db.Exec(
		"INSERT INTO projects (title, description, status) VALUES (?, ?, ?)",
		"Invalid Project", "Description", "pending",
	)
	require.Error(t, err, "should error on invalid project status")
}

func TestSchema_ForeignKeyProjectMembers(t *testing.T) {
	db := helpers.NewTestDB(t)

	// Create test data
	projectID := helpers.InsertProject(t, db, "Test Project", "Description", "active")
	memberID := helpers.InsertLabMember(t, db, "Test Member", "PhD")

	// Test valid insert
	_, err := db.Exec(
		"INSERT INTO project_members (project_id, member_id) VALUES (?, ?)",
		projectID, memberID,
	)
	require.NoError(t, err)

	// Test invalid project_id
	_, err = db.Exec(
		"INSERT INTO project_members (project_id, member_id) VALUES (?, ?)",
		99999, memberID,
	)
	require.Error(t, err, "should error on invalid project_id")
	require.Contains(t, err.Error(), "FOREIGN KEY")

	// Test invalid member_id
	_, err = db.Exec(
		"INSERT INTO project_members (project_id, member_id) VALUES (?, ?)",
		projectID, 99999,
	)
	require.Error(t, err, "should error on invalid member_id")
	require.Contains(t, err.Error(), "FOREIGN KEY")

	// Test cascade delete
	_, err = db.Exec("DELETE FROM projects WHERE id = ?", projectID)
	require.NoError(t, err)

	var count int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM project_members WHERE project_id = ?",
		projectID,
	).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count, "junction record should be deleted on cascade")
}

func TestSchema_ForeignKeyPublicationAuthors(t *testing.T) {
	db := helpers.NewTestDB(t)

	// Create test data
	pubID := helpers.InsertPublication(t, db, "Test Publication", "Author Name", 2024)
	memberID := helpers.InsertLabMember(t, db, "Test Author", "PhD")

	// Test valid insert
	_, err := db.Exec(
		"INSERT INTO publication_authors (publication_id, member_id) VALUES (?, ?)",
		pubID, memberID,
	)
	require.NoError(t, err)

	// Test cascade delete on publication
	_, err = db.Exec("DELETE FROM publications WHERE id = ?", pubID)
	require.NoError(t, err)

	var count int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM publication_authors WHERE publication_id = ?",
		pubID,
	).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count, "junction record should be deleted on cascade")
}

func TestSchema_ForeignKeyProjectPublications(t *testing.T) {
	db := helpers.NewTestDB(t)

	// Create test data
	projectID := helpers.InsertProject(t, db, "Test Project", "Description", "active")
	pubID := helpers.InsertPublication(t, db, "Test Publication", "Author", 2024)

	// Test valid insert
	_, err := db.Exec(
		"INSERT INTO project_publications (project_id, publication_id) VALUES (?, ?)",
		projectID, pubID,
	)
	require.NoError(t, err)

	// Test cascade delete on project
	_, err = db.Exec("DELETE FROM projects WHERE id = ?", projectID)
	require.NoError(t, err)

	var count int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM project_publications WHERE project_id = ?",
		projectID,
	).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count, "junction record should be deleted on cascade")
}

func TestSchema_HomepageSectionKeyUnique(t *testing.T) {
	db := helpers.NewTestDB(t)

	// Insert first section
	_, err := db.Exec(
		"INSERT INTO homepage_sections (section_key, title, content) VALUES (?, ?, ?)",
		"overview", "Overview", "Content",
	)
	require.NoError(t, err)

	// Attempt duplicate section_key
	_, err = db.Exec(
		"INSERT INTO homepage_sections (section_key, title, content) VALUES (?, ?, ?)",
		"overview", "Duplicate", "Content",
	)
	require.Error(t, err, "should error on duplicate section_key")
	require.Contains(t, err.Error(), "UNIQUE")
}

func TestSchema_NotNullConstraints(t *testing.T) {
	db := helpers.NewTestDB(t)

	// Test users table NOT NULL constraints
	_, err := db.Exec("INSERT INTO users (email) VALUES (?)", "test@example.com")
	require.Error(t, err, "should error when password_hash is NULL")

	_, err = db.Exec("INSERT INTO users (password_hash) VALUES (?)", "hash")
	require.Error(t, err, "should error when email is NULL")

	// Test lab_members table NOT NULL constraints
	_, err = db.Exec("INSERT INTO lab_members (role) VALUES (?)", "PhD")
	require.Error(t, err, "should error when name is NULL")

	_, err = db.Exec("INSERT INTO lab_members (name) VALUES (?)", "Test Member")
	require.Error(t, err, "should error when role is NULL")
}
