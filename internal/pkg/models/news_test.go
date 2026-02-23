package models

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNews_Validation(t *testing.T) {
	v := newValidator()

	validNews := News{
		Title:   "Test News Title",
		Content: "This is the content of the news item.",
	}

	err := validateStruct(v, validNews)
	assert.NoError(t, err, "valid news should pass validation")
}

func TestNews_Validation_MissingTitle(t *testing.T) {
	v := newValidator()

	news := News{
		Title:   "",
		Content: "Content here",
	}

	err := validateStruct(v, news)
	assert.Error(t, err, "empty title should fail validation")
}

func TestNews_Validation_MissingContent(t *testing.T) {
	v := newValidator()

	news := News{
		Title:   "Test Title",
		Content: "",
	}

	err := validateStruct(v, news)
	assert.Error(t, err, "empty content should fail validation")
}

func TestNews_IsPublishedNow_Draft(t *testing.T) {
	news := News{
		ID:          1,
		Title:       "Draft News",
		Content:     "This is a draft",
		IsPublished: false,
		PublishedAt: sql.NullTime{Valid: false},
	}

	assert.False(t, news.IsPublishedNow(), "draft news should not be published")
}

func TestNews_IsPublishedNow_Future(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)
	news := News{
		ID:          2,
		Title:       "Scheduled News",
		Content:     "This is scheduled for future",
		IsPublished: true,
		PublishedAt: sql.NullTime{
			Time:  futureTime,
			Valid: true,
		},
	}

	assert.False(t, news.IsPublishedNow(), "future news should not be published yet")
}

func TestNews_IsPublishedNow_Current(t *testing.T) {
	now := time.Now()
	news := News{
		ID:          3,
		Title:       "Current News",
		Content:     "This is published now",
		IsPublished: true,
		PublishedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	assert.True(t, news.IsPublishedNow(), "current news should be published")
}

func TestNews_IsPublishedNow_Past(t *testing.T) {
	pastTime := time.Now().Add(-24 * time.Hour)
	news := News{
		ID:          4,
		Title:       "Past News",
		Content:     "This was published yesterday",
		IsPublished: true,
		PublishedAt: sql.NullTime{
			Time:  pastTime,
			Valid: true,
		},
	}

	assert.True(t, news.IsPublishedNow(), "past news should be published")
}

func TestNews_IsPublishedNow_NullPublishedAt(t *testing.T) {
	news := News{
		ID:          5,
		Title:       "News with no date",
		Content:     "This has no publish date",
		IsPublished: true,
		PublishedAt: sql.NullTime{Valid: false},
	}

	assert.False(t, news.IsPublishedNow(), "news with null PublishedAt should not be published")
}

func TestNews_JSONSerialization(t *testing.T) {
	publishedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	news := News{
		ID:          1,
		Title:       "Test News",
		Content:     "News content here",
		PublishedAt: sql.NullTime{Time: publishedTime, Valid: true},
		IsPublished: true,
	}

	data, err := json.Marshal(news)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"id\":1")
	assert.Contains(t, jsonStr, "\"title\":\"Test News\"")
	assert.Contains(t, jsonStr, "\"content\":\"News content here\"")
	assert.Contains(t, jsonStr, "\"is_published\":true")
	assert.Contains(t, jsonStr, "\"published_at\"")
}

func TestNews_JSON_NullableOmit(t *testing.T) {
	news := News{
		ID:          2,
		Title:       "Draft News",
		Content:     "Draft content",
		PublishedAt: sql.NullTime{Valid: false},
		IsPublished: false,
	}

	data, err := json.Marshal(news)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "\"title\":\"Draft News\"")
	// sql.NullTime fields are not omitted, they appear as {"Time":"0001-01-01...", "Valid":false}
	assert.Contains(t, jsonStr, "\"published_at\"")
	assert.Contains(t, jsonStr, "\"Valid\":false")
}

func TestNews_NullableTime(t *testing.T) {
	// Test valid time
	validTime := sql.NullTime{
		Time:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Valid: true,
	}
	assert.True(t, validTime.Valid)
	assert.Equal(t, 2024, validTime.Time.Year())

	// Test null time
	nullTime := sql.NullTime{
		Valid: false,
	}
	assert.False(t, nullTime.Valid)
}

func TestNews_JSONDeserialization(t *testing.T) {
	// sql.NullTime fields need to be deserialized with proper object format
	jsonData := `{
		"id": 1,
		"title": "Test News",
		"content": "Test Content",
		"is_published": true,
		"published_at": {"Time": "2024-01-15T10:30:00Z", "Valid": true}
	}`

	var news News
	err := json.Unmarshal([]byte(jsonData), &news)
	require.NoError(t, err)

	assert.Equal(t, 1, news.ID)
	assert.Equal(t, "Test News", news.Title)
	assert.Equal(t, "Test Content", news.Content)
	assert.True(t, news.IsPublished)
	assert.True(t, news.PublishedAt.Valid)
}
