package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProject_Validation(t *testing.T) {
	v := newValidator()

	validProject := Project{
		Title:       "Test Project",
		Description: "A test project description",
		Status:      ProjectStatusActive,
	}

	err := validateStruct(v, validProject)
	assert.NoError(t, err, "valid project should pass validation")
}

func TestProject_Validation_MissingTitle(t *testing.T) {
	v := newValidator()

	project := Project{
		Title:       "",
		Description: "Description",
		Status:      ProjectStatusActive,
	}

	err := validateStruct(v, project)
	assert.Error(t, err, "empty title should fail validation")
}

func TestProject_Validation_MissingDescription(t *testing.T) {
	v := newValidator()

	project := Project{
		Title:       "Test Project",
		Description: "",
		Status:      ProjectStatusActive,
	}

	err := validateStruct(v, project)
	assert.Error(t, err, "empty description should fail validation")
}

func TestProject_Validation_InvalidStatus(t *testing.T) {
	v := newValidator()

	tests := []struct {
		name   string
		status string
	}{
		{"empty status", ""},
		{"invalid status", "pending"},
		{"wrong case", "Active"},
		{"typo", "complet"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := Project{
				Title:       "Test Project",
				Description: "Description",
				Status:      ProjectStatus(tt.status),
			}
			err := validateStruct(v, project)
			assert.Error(t, err, "invalid status should fail validation")
		})
	}
}

func TestProject_Validation_ValidStatuses(t *testing.T) {
	v := newValidator()

	statuses := []ProjectStatus{ProjectStatusActive, ProjectStatusCompleted}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			project := Project{
				Title:       "Test Project",
				Description: "Description",
				Status:      status,
			}
			err := validateStruct(v, project)
			assert.NoError(t, err, "status %s should be valid", status)
		})
	}
}

func TestProject_JSONSerialization(t *testing.T) {
	project := Project{
		ID:          1,
		Title:       "AI Research Project",
		Description: "Researching machine learning applications",
		Status:      ProjectStatusActive,
	}

	data, err := json.Marshal(project)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"id\":1")
	assert.Contains(t, jsonStr, "\"title\":\"AI Research Project\"")
	assert.Contains(t, jsonStr, "\"description\":\"Researching machine learning applications\"")
	assert.Contains(t, jsonStr, "\"status\":\"active\"")
}

func TestProject_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"id": 1,
		"title": "Test Project",
		"description": "Test Description",
		"status": "completed"
	}`

	var project Project
	err := json.Unmarshal([]byte(jsonData), &project)
	require.NoError(t, err)

	assert.Equal(t, 1, project.ID)
	assert.Equal(t, "Test Project", project.Title)
	assert.Equal(t, "Test Description", project.Description)
	assert.Equal(t, ProjectStatusCompleted, project.Status)
}
