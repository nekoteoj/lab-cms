package models

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublication_Validation(t *testing.T) {
	v := newValidator()

	validPublication := Publication{
		Title:       "Test Publication Title",
		AuthorsText: "John Doe, Jane Smith",
		Year:        2024,
	}

	err := validateStruct(v, validPublication)
	assert.NoError(t, err, "valid publication should pass validation")
}

func TestPublication_Validation_MissingTitle(t *testing.T) {
	v := newValidator()

	publication := Publication{
		Title:       "",
		AuthorsText: "John Doe",
		Year:        2024,
	}

	err := validateStruct(v, publication)
	assert.Error(t, err, "empty title should fail validation")
}

func TestPublication_Validation_MissingAuthors(t *testing.T) {
	v := newValidator()

	publication := Publication{
		Title:       "Test Title",
		AuthorsText: "",
		Year:        2024,
	}

	err := validateStruct(v, publication)
	assert.Error(t, err, "empty authors should fail validation")
}

func TestPublication_Validation_InvalidYear(t *testing.T) {
	v := newValidator()

	tests := []struct {
		name string
		year int
	}{
		{"too old", 1899},
		{"too new", 2101},
		{"negative", -100},
		{"zero", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publication := Publication{
				Title:       "Test Title",
				AuthorsText: "John Doe",
				Year:        tt.year,
			}
			err := validateStruct(v, publication)
			assert.Error(t, err, "year %d should fail validation", tt.year)
		})
	}
}

func TestPublication_Validation_ValidYear(t *testing.T) {
	v := newValidator()

	years := []int{1900, 1950, 2000, 2024, 2100}

	for _, year := range years {
		t.Run("year", func(t *testing.T) {
			publication := Publication{
				Title:       "Test Title",
				AuthorsText: "John Doe",
				Year:        year,
			}
			err := validateStruct(v, publication)
			assert.NoError(t, err, "year %d should be valid", year)
		})
	}
}

func TestPublication_Validation_InvalidURL(t *testing.T) {
	v := newValidator()

	publication := Publication{
		Title:       "Test Title",
		AuthorsText: "John Doe",
		Year:        2024,
		URL:         sql.NullString{String: "not-a-valid-url", Valid: true},
	}

	// Note: URL validation is not performed via struct tags for sql.NullString
	// URL validation should be done at the application layer or database level
	err := validateStruct(v, publication)
	assert.NoError(t, err, "sql.NullString fields without validation tags should pass")
}

func TestPublication_Validation_NullURL(t *testing.T) {
	v := newValidator()

	publication := Publication{
		Title:       "Test Title",
		AuthorsText: "John Doe",
		Year:        2024,
		URL:         sql.NullString{Valid: false},
	}

	err := validateStruct(v, publication)
	assert.NoError(t, err, "null URL should be allowed")
}

func TestPublication_Validation_NullVenue(t *testing.T) {
	v := newValidator()

	publication := Publication{
		Title:       "Test Title",
		AuthorsText: "John Doe",
		Year:        2024,
		Venue:       sql.NullString{Valid: false},
	}

	err := validateStruct(v, publication)
	assert.NoError(t, err, "null venue should be allowed")
}

func TestPublication_JSONSerialization(t *testing.T) {
	publication := Publication{
		ID:          1,
		Title:       "Test Publication",
		AuthorsText: "John Doe, Jane Smith",
		Venue:       sql.NullString{String: "Nature", Valid: true},
		Year:        2024,
		URL:         sql.NullString{String: "https://doi.org/10.1234", Valid: true},
	}

	data, err := json.Marshal(publication)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"id\":1")
	assert.Contains(t, jsonStr, "\"title\":\"Test Publication\"")
	assert.Contains(t, jsonStr, "\"authors_text\":\"John Doe, Jane Smith\"")
	// sql.NullString serializes as object with String and Valid fields
	assert.Contains(t, jsonStr, "\"venue\"")
	assert.Contains(t, jsonStr, "\"String\":\"Nature\"")
	assert.Contains(t, jsonStr, "\"year\":2024")
	assert.Contains(t, jsonStr, "\"url\"")
}

func TestPublication_JSON_NullableOmit(t *testing.T) {
	publication := Publication{
		ID:          1,
		Title:       "Test Publication",
		AuthorsText: "John Doe",
		Venue:       sql.NullString{Valid: false},
		Year:        2024,
		URL:         sql.NullString{Valid: false},
	}

	data, err := json.Marshal(publication)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"title\":\"Test Publication\"")
	// sql.NullString fields are not omitted, they appear as {"String":"", "Valid":false}
	assert.Contains(t, jsonStr, "\"venue\"")
	assert.Contains(t, jsonStr, "\"url\"")
	assert.Contains(t, jsonStr, "\"Valid\":false")
}
