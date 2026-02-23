package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicationWithAuthors(t *testing.T) {
	publication := Publication{
		ID:          1,
		Title:       "Test Publication",
		AuthorsText: "John Doe, Jane Smith",
		Year:        2024,
	}

	authors := []LabMember{
		{
			ID:   1,
			Name: "John Doe",
			Role: LabMemberRolePhD,
		},
		{
			ID:   2,
			Name: "Jane Smith",
			Role: LabMemberRolePostdoc,
		},
	}

	pubWithAuthors := PublicationWithAuthors{
		Publication: publication,
		Authors:     authors,
	}

	assert.Equal(t, publication.ID, pubWithAuthors.ID)
	assert.Equal(t, publication.Title, pubWithAuthors.Title)
	assert.Len(t, pubWithAuthors.Authors, 2)
	assert.Equal(t, "John Doe", pubWithAuthors.Authors[0].Name)
	assert.Equal(t, "Jane Smith", pubWithAuthors.Authors[1].Name)
}

func TestPublicationWithAuthors_JSONSerialization(t *testing.T) {
	publication := Publication{
		ID:          1,
		Title:       "Test Publication",
		AuthorsText: "John Doe, Jane Smith",
		Year:        2024,
	}

	authors := []LabMember{
		{
			ID:   1,
			Name: "John Doe",
			Role: LabMemberRolePhD,
		},
	}

	pubWithAuthors := PublicationWithAuthors{
		Publication: publication,
		Authors:     authors,
	}

	data, err := json.Marshal(pubWithAuthors)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"id\":1")
	assert.Contains(t, jsonStr, "\"title\":\"Test Publication\"")
	assert.Contains(t, jsonStr, "\"authors\"")
	assert.Contains(t, jsonStr, "\"John Doe\"")
	assert.Contains(t, jsonStr, "\"PhD\"")
}

func TestPublicationWithAuthors_JSON_NilAuthors(t *testing.T) {
	publication := Publication{
		ID:          1,
		Title:       "Test Publication",
		AuthorsText: "John Doe",
		Year:        2024,
	}

	pubWithAuthors := PublicationWithAuthors{
		Publication: publication,
		Authors:     nil,
	}

	data, err := json.Marshal(pubWithAuthors)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"title\":\"Test Publication\"")
	// Authors field should appear as null when nil since omitempty was removed
	assert.Contains(t, jsonStr, "\"authors\":null")
}

func TestProjectWithRelations(t *testing.T) {
	project := Project{
		ID:          1,
		Title:       "AI Research Project",
		Description: "Machine learning research",
		Status:      ProjectStatusActive,
	}

	members := []LabMember{
		{
			ID:   1,
			Name: "John Doe",
			Role: LabMemberRolePhD,
		},
	}

	publications := []Publication{
		{
			ID:          1,
			Title:       "Related Paper",
			AuthorsText: "John Doe",
			Year:        2024,
		},
	}

	projWithRelations := ProjectWithRelations{
		Project:      project,
		Members:      members,
		Publications: publications,
	}

	assert.Equal(t, project.ID, projWithRelations.ID)
	assert.Equal(t, project.Title, projWithRelations.Title)
	assert.Len(t, projWithRelations.Members, 1)
	assert.Len(t, projWithRelations.Publications, 1)
	assert.Equal(t, "John Doe", projWithRelations.Members[0].Name)
	assert.Equal(t, "Related Paper", projWithRelations.Publications[0].Title)
}

func TestProjectWithRelations_JSONSerialization(t *testing.T) {
	project := Project{
		ID:          1,
		Title:       "AI Research Project",
		Description: "Machine learning research",
		Status:      ProjectStatusActive,
		CreatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	members := []LabMember{
		{
			ID:   1,
			Name: "John Doe",
			Role: LabMemberRolePhD,
		},
	}

	publications := []Publication{
		{
			ID:          1,
			Title:       "Related Paper",
			AuthorsText: "John Doe",
			Year:        2024,
		},
	}

	projWithRelations := ProjectWithRelations{
		Project:      project,
		Members:      members,
		Publications: publications,
	}

	data, err := json.Marshal(projWithRelations)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"id\":1")
	assert.Contains(t, jsonStr, "\"title\":\"AI Research Project\"")
	assert.Contains(t, jsonStr, "\"members\"")
	assert.Contains(t, jsonStr, "\"publications\"")
	assert.Contains(t, jsonStr, "\"John Doe\"")
	assert.Contains(t, jsonStr, "\"Related Paper\"")
}

func TestProjectWithRelations_JSON_EmptyRelations(t *testing.T) {
	project := Project{
		ID:          1,
		Title:       "AI Research Project",
		Description: "Machine learning research",
		Status:      ProjectStatusActive,
	}

	projWithRelations := ProjectWithRelations{
		Project:      project,
		Members:      []LabMember{},
		Publications: []Publication{},
	}

	data, err := json.Marshal(projWithRelations)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"title\":\"AI Research Project\"")
	assert.Contains(t, jsonStr, "\"members\":[]")
	assert.Contains(t, jsonStr, "\"publications\":[]")
}

func TestExtendedStructs_Embedding(t *testing.T) {
	// Test that embedded structs are properly accessible
	publication := Publication{
		ID:    1,
		Title: "Test",
		Year:  2024,
	}

	pubWithAuthors := PublicationWithAuthors{
		Publication: publication,
		Authors:     nil,
	}

	// Should be able to access embedded fields directly
	assert.Equal(t, 1, pubWithAuthors.Publication.ID)
	assert.Equal(t, "Test", pubWithAuthors.Publication.Title)

	project := Project{
		ID:     2,
		Title:  "Test Project",
		Status: ProjectStatusCompleted,
	}

	projWithRelations := ProjectWithRelations{
		Project: project,
	}

	assert.Equal(t, 2, projWithRelations.Project.ID)
	assert.Equal(t, "Test Project", projWithRelations.Project.Title)
	assert.Equal(t, ProjectStatusCompleted, projWithRelations.Project.Status)
}
