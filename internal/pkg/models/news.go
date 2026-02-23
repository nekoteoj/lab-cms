package models

import (
	"database/sql"
	"time"
)

// News represents a news item or announcement
type News struct {
	ID          int          `json:"id"`
	Title       string       `json:"title" validate:"required,max=255"`
	Content     string       `json:"content" validate:"required"`
	PublishedAt sql.NullTime `json:"published_at,omitempty"`
	IsPublished bool         `json:"is_published"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// IsPublishedNow returns true if the news item should be visible to the public
func (n *News) IsPublishedNow() bool {
	if !n.IsPublished {
		return false
	}
	if !n.PublishedAt.Valid {
		return false
	}
	return n.PublishedAt.Time.Before(time.Now()) || n.PublishedAt.Time.Equal(time.Now())
}
