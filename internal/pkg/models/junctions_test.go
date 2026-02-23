package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectMember_Validation(t *testing.T) {
	v := newValidator()

	validJunction := ProjectMember{
		ProjectID: 1,
		MemberID:  2,
	}

	err := validateStruct(v, validJunction)
	assert.NoError(t, err, "valid project member should pass validation")
}

func TestProjectMember_Validation_MissingProjectID(t *testing.T) {
	v := newValidator()

	junction := ProjectMember{
		ProjectID: 0,
		MemberID:  2,
	}

	err := validateStruct(v, junction)
	assert.Error(t, err, "missing project ID should fail validation")
}

func TestProjectMember_Validation_MissingMemberID(t *testing.T) {
	v := newValidator()

	junction := ProjectMember{
		ProjectID: 1,
		MemberID:  0,
	}

	err := validateStruct(v, junction)
	assert.Error(t, err, "missing member ID should fail validation")
}

func TestProjectMember_JSONSerialization(t *testing.T) {
	junction := ProjectMember{
		ProjectID: 1,
		MemberID:  2,
	}

	data, err := json.Marshal(junction)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"project_id\":1")
	assert.Contains(t, jsonStr, "\"member_id\":2")
}

func TestPublicationAuthor_Validation(t *testing.T) {
	v := newValidator()

	validJunction := PublicationAuthor{
		PublicationID: 1,
		MemberID:      2,
	}

	err := validateStruct(v, validJunction)
	assert.NoError(t, err, "valid publication author should pass validation")
}

func TestPublicationAuthor_Validation_MissingIDs(t *testing.T) {
	v := newValidator()

	junction := PublicationAuthor{
		PublicationID: 0,
		MemberID:      0,
	}

	err := validateStruct(v, junction)
	assert.Error(t, err, "missing IDs should fail validation")
}

func TestPublicationAuthor_JSONSerialization(t *testing.T) {
	junction := PublicationAuthor{
		PublicationID: 1,
		MemberID:      2,
	}

	data, err := json.Marshal(junction)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"publication_id\":1")
	assert.Contains(t, jsonStr, "\"member_id\":2")
}

func TestProjectPublication_Validation(t *testing.T) {
	v := newValidator()

	validJunction := ProjectPublication{
		ProjectID:     1,
		PublicationID: 2,
	}

	err := validateStruct(v, validJunction)
	assert.NoError(t, err, "valid project publication should pass validation")
}

func TestProjectPublication_Validation_MissingIDs(t *testing.T) {
	v := newValidator()

	junction := ProjectPublication{
		ProjectID:     0,
		PublicationID: 0,
	}

	err := validateStruct(v, junction)
	assert.Error(t, err, "missing IDs should fail validation")
}

func TestProjectPublication_JSONSerialization(t *testing.T) {
	junction := ProjectPublication{
		ProjectID:     1,
		PublicationID: 2,
	}

	data, err := json.Marshal(junction)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"project_id\":1")
	assert.Contains(t, jsonStr, "\"publication_id\":2")
}

func TestJunctionStructs_JSONDeserialization(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		target   interface{}
		expected interface{}
	}{
		{
			name:   "ProjectMember",
			json:   `{"project_id":1,"member_id":2}`,
			target: &ProjectMember{},
			expected: &ProjectMember{
				ProjectID: 1,
				MemberID:  2,
			},
		},
		{
			name:   "PublicationAuthor",
			json:   `{"publication_id":3,"member_id":4}`,
			target: &PublicationAuthor{},
			expected: &PublicationAuthor{
				PublicationID: 3,
				MemberID:      4,
			},
		},
		{
			name:   "ProjectPublication",
			json:   `{"project_id":5,"publication_id":6}`,
			target: &ProjectPublication{},
			expected: &ProjectPublication{
				ProjectID:     5,
				PublicationID: 6,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tt.json), tt.target)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, tt.target)
		})
	}
}
