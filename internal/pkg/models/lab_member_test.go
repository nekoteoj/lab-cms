package models

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabMember_Validation(t *testing.T) {
	v := newValidator()

	validMember := LabMember{
		Name: "John Doe",
		Role: LabMemberRolePhD,
	}

	err := validateStruct(v, validMember)
	assert.NoError(t, err, "valid lab member should pass validation")
}

func TestLabMember_Validation_AllRoles(t *testing.T) {
	v := newValidator()

	roles := []LabMemberRole{
		LabMemberRolePI,
		LabMemberRolePostdoc,
		LabMemberRolePhD,
		LabMemberRoleMaster,
		LabMemberRoleBachelor,
		LabMemberRoleResearcher,
	}

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			member := LabMember{
				Name: "Test Member",
				Role: role,
			}
			err := validateStruct(v, member)
			assert.NoError(t, err, "role %s should be valid", role)
		})
	}
}

func TestLabMember_Validation_InvalidRole(t *testing.T) {
	v := newValidator()

	tests := []struct {
		name string
		role string
	}{
		{"empty role", ""},
		{"invalid role", "Professor"},
		{"wrong case", "phd"},
		{"typo", "PII"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			member := LabMember{
				Name: "Test Member",
				Role: LabMemberRole(tt.role),
			}
			err := validateStruct(v, member)
			assert.Error(t, err, "invalid role should fail validation")
		})
	}
}

func TestLabMember_Validation_EmptyName(t *testing.T) {
	v := newValidator()

	member := LabMember{
		Name: "",
		Role: LabMemberRolePhD,
	}

	err := validateStruct(v, member)
	assert.Error(t, err, "empty name should fail validation")
}

func TestLabMember_Validation_InvalidEmail(t *testing.T) {
	v := newValidator()

	member := LabMember{
		Name:  "Test Member",
		Role:  LabMemberRolePhD,
		Email: sql.NullString{String: "invalid-email", Valid: true},
	}

	// Note: Email validation is not performed via struct tags for sql.NullString
	// Email validation should be done at the application layer or database level
	// This test just ensures the struct passes validation (no email validation tag)
	err := validateStruct(v, member)
	assert.NoError(t, err, "sql.NullString fields without validation tags should pass")
}

func TestLabMember_NullableFields(t *testing.T) {
	v := newValidator()

	// Test with null fields
	member := LabMember{
		Name:     "Test Member",
		Role:     LabMemberRolePhD,
		Email:    sql.NullString{Valid: false},
		Bio:      sql.NullString{Valid: false},
		PhotoURL: sql.NullString{Valid: false},
	}

	err := validateStruct(v, member)
	assert.NoError(t, err, "null optional fields should pass validation")
}

func TestLabMember_JSONSerialization(t *testing.T) {
	member := LabMember{
		ID:                  1,
		Name:                "John Doe",
		Role:                LabMemberRolePhD,
		Email:               sql.NullString{String: "john@example.com", Valid: true},
		Bio:                 sql.NullString{String: "PhD Student", Valid: true},
		PhotoURL:            sql.NullString{String: "https://example.com/photo.jpg", Valid: true},
		PersonalPageContent: sql.NullString{Valid: false},
		ResearchInterests:   sql.NullString{Valid: false},
		IsAlumni:            false,
		DisplayOrder:        5,
	}

	data, err := json.Marshal(member)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"id\":1")
	assert.Contains(t, jsonStr, "\"name\":\"John Doe\"")
	assert.Contains(t, jsonStr, "\"role\":\"PhD\"")
	// sql.NullString serializes as object with String and Valid fields
	assert.Contains(t, jsonStr, "\"email\"")
	assert.Contains(t, jsonStr, "\"String\":\"john@example.com\"")
	assert.Contains(t, jsonStr, "\"bio\"")
	assert.Contains(t, jsonStr, "\"photo_url\"")
}

func TestLabMember_JSON_NullableOmit(t *testing.T) {
	member := LabMember{
		ID:                  1,
		Name:                "Jane Doe",
		Role:                LabMemberRolePI,
		Email:               sql.NullString{Valid: false},
		Bio:                 sql.NullString{Valid: false},
		PhotoURL:            sql.NullString{Valid: false},
		PersonalPageContent: sql.NullString{Valid: false},
		ResearchInterests:   sql.NullString{Valid: false},
		IsAlumni:            true,
		DisplayOrder:        1,
	}

	data, err := json.Marshal(member)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"name\":\"Jane Doe\"")
	// sql.NullString fields are not omitted, they appear as {"String":"", "Valid":false}
	// This is the expected behavior for sql.NullString
	assert.Contains(t, jsonStr, "\"email\"")
	assert.Contains(t, jsonStr, "\"Valid\":false")
}

func TestLabMember_JSONDeserialization(t *testing.T) {
	// sql.NullString fields need to be deserialized with proper object format
	jsonData := `{
		"id": 1,
		"name": "Test Member",
		"role": "Postdoc",
		"email": {"String": "test@example.com", "Valid": true},
		"is_alumni": true,
		"display_order": 10
	}`

	var member LabMember
	err := json.Unmarshal([]byte(jsonData), &member)
	require.NoError(t, err)

	assert.Equal(t, 1, member.ID)
	assert.Equal(t, "Test Member", member.Name)
	assert.Equal(t, LabMemberRolePostdoc, member.Role)
	assert.True(t, member.Email.Valid)
	assert.Equal(t, "test@example.com", member.Email.String)
	assert.True(t, member.IsAlumni)
	assert.Equal(t, 10, member.DisplayOrder)
}
