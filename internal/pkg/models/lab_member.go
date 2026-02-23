package models

import (
	"database/sql"
	"time"
)

// LabMember represents a lab member (PI, Postdoc, PhD, etc.)
type LabMember struct {
	ID                  int            `json:"id"`
	Name                string         `json:"name" validate:"required,max=255"`
	Role                LabMemberRole  `json:"role" validate:"required,oneof=PI Postdoc PhD Master Bachelor Researcher"`
	Email               sql.NullString `json:"email,omitempty"`
	Bio                 sql.NullString `json:"bio,omitempty"`
	PhotoURL            sql.NullString `json:"photo_url,omitempty"`
	PersonalPageContent sql.NullString `json:"personal_page_content,omitempty"`
	ResearchInterests   sql.NullString `json:"research_interests,omitempty"`
	IsAlumni            bool           `json:"is_alumni"`
	DisplayOrder        int            `json:"display_order"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}
