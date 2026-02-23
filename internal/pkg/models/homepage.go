package models

import (
	"time"
)

// HomepageSection represents an editable section of the homepage
type HomepageSection struct {
	ID           int       `json:"id"`
	SectionKey   string    `json:"section_key" validate:"required,max=100"`
	Title        string    `json:"title" validate:"required,max=255"`
	Content      string    `json:"content" validate:"required"`
	DisplayOrder int       `json:"display_order"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Common section keys for the homepage
const (
	HomepageSectionOverview = "overview"
	HomepageSectionMission  = "mission"
	HomepageSectionResearch = "research"
	HomepageSectionContact  = "contact"
)
