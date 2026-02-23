package models

import (
	"database/sql"
	"time"
)

// Publication represents a research publication
type Publication struct {
	ID          int            `json:"id"`
	Title       string         `json:"title" validate:"required,max=500"`
	AuthorsText string         `json:"authors_text" validate:"required"`
	Venue       sql.NullString `json:"venue,omitempty"`
	Year        int            `json:"year" validate:"required,min=1900,max=2100"`
	URL         sql.NullString `json:"url,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// PublicationWithAuthors extends Publication to include associated lab members
type PublicationWithAuthors struct {
	Publication
	Authors []LabMember `json:"authors"`
}
