package models

// ProjectMember represents the many-to-many relationship between projects and lab members
type ProjectMember struct {
	ProjectID int `json:"project_id" validate:"required"`
	MemberID  int `json:"member_id" validate:"required"`
}

// PublicationAuthor represents the many-to-many relationship between publications and lab members
type PublicationAuthor struct {
	PublicationID int `json:"publication_id" validate:"required"`
	MemberID      int `json:"member_id" validate:"required"`
}

// ProjectPublication represents the many-to-many relationship between projects and publications
type ProjectPublication struct {
	ProjectID     int `json:"project_id" validate:"required"`
	PublicationID int `json:"publication_id" validate:"required"`
}
