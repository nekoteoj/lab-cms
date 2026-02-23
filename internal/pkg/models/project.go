package models

import (
	"time"
)

// Project represents a research project
type Project struct {
	ID          int           `json:"id"`
	Title       string        `json:"title" validate:"required,max=255"`
	Description string        `json:"description" validate:"required"`
	Status      ProjectStatus `json:"status" validate:"required,oneof=active completed"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// ProjectWithRelations extends Project to include associated members and publications
type ProjectWithRelations struct {
	Project
	Members      []LabMember   `json:"members"`
	Publications []Publication `json:"publications"`
}
