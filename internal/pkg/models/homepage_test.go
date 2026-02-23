package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHomepageSection_Validation(t *testing.T) {
	v := newValidator()

	validSection := HomepageSection{
		SectionKey:   "overview",
		Title:        "Lab Overview",
		Content:      "Welcome to our lab",
		DisplayOrder: 1,
	}

	err := validateStruct(v, validSection)
	assert.NoError(t, err, "valid homepage section should pass validation")
}

func TestHomepageSection_Validation_MissingSectionKey(t *testing.T) {
	v := newValidator()

	section := HomepageSection{
		SectionKey:   "",
		Title:        "Lab Overview",
		Content:      "Welcome to our lab",
		DisplayOrder: 1,
	}

	err := validateStruct(v, section)
	assert.Error(t, err, "empty section key should fail validation")
}

func TestHomepageSection_Validation_MissingTitle(t *testing.T) {
	v := newValidator()

	section := HomepageSection{
		SectionKey:   "overview",
		Title:        "",
		Content:      "Welcome to our lab",
		DisplayOrder: 1,
	}

	err := validateStruct(v, section)
	assert.Error(t, err, "empty title should fail validation")
}

func TestHomepageSection_Validation_MissingContent(t *testing.T) {
	v := newValidator()

	section := HomepageSection{
		SectionKey:   "overview",
		Title:        "Lab Overview",
		Content:      "",
		DisplayOrder: 1,
	}

	err := validateStruct(v, section)
	assert.Error(t, err, "empty content should fail validation")
}

func TestHomepageSection_Validation_TooLongSectionKey(t *testing.T) {
	v := newValidator()

	// Create a string that is definitely longer than 100 characters (this one is 104)
	longKey := "this-is-a-very-long-section-key-that-definitely-exceeds-one-hundred-characters-limit-and-should-fail-validation"

	section := HomepageSection{
		SectionKey:   longKey,
		Title:        "Title",
		Content:      "Content",
		DisplayOrder: 1,
	}

	err := validateStruct(v, section)
	assert.Error(t, err, "section key longer than 100 characters should fail validation")
}

func TestHomepageSection_Validation_TooLongTitle(t *testing.T) {
	v := newValidator()

	longTitle := make([]byte, 256)
	for i := range longTitle {
		longTitle[i] = 'a'
	}

	section := HomepageSection{
		SectionKey:   "overview",
		Title:        string(longTitle),
		Content:      "Content",
		DisplayOrder: 1,
	}

	err := validateStruct(v, section)
	assert.Error(t, err, "title longer than 255 characters should fail validation")
}

func TestHomepageSectionConstants(t *testing.T) {
	assert.Equal(t, "overview", HomepageSectionOverview)
	assert.Equal(t, "mission", HomepageSectionMission)
	assert.Equal(t, "research", HomepageSectionResearch)
	assert.Equal(t, "contact", HomepageSectionContact)
}

func TestHomepageSection_JSONSerialization(t *testing.T) {
	section := HomepageSection{
		ID:           1,
		SectionKey:   "overview",
		Title:        "Lab Overview",
		Content:      "Welcome to our research lab",
		DisplayOrder: 5,
	}

	data, err := json.Marshal(section)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"id\":1")
	assert.Contains(t, jsonStr, "\"section_key\":\"overview\"")
	assert.Contains(t, jsonStr, "\"title\":\"Lab Overview\"")
	assert.Contains(t, jsonStr, "\"content\":\"Welcome to our research lab\"")
	assert.Contains(t, jsonStr, "\"display_order\":5")
}

func TestHomepageSection_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"id": 1,
		"section_key": "mission",
		"title": "Our Mission",
		"content": "To advance science",
		"display_order": 2
	}`

	var section HomepageSection
	err := json.Unmarshal([]byte(jsonData), &section)
	require.NoError(t, err)

	assert.Equal(t, 1, section.ID)
	assert.Equal(t, "mission", section.SectionKey)
	assert.Equal(t, "Our Mission", section.Title)
	assert.Equal(t, "To advance science", section.Content)
	assert.Equal(t, 2, section.DisplayOrder)
}
