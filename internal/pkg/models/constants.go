package models

// UserRole defines the possible roles for users
type UserRole string

const (
	UserRoleNormal UserRole = "normal"
	UserRoleRoot   UserRole = "root"
)

// LabMemberRole defines the possible roles for lab members
type LabMemberRole string

const (
	LabMemberRolePI         LabMemberRole = "PI"
	LabMemberRolePostdoc    LabMemberRole = "Postdoc"
	LabMemberRolePhD        LabMemberRole = "PhD"
	LabMemberRoleMaster     LabMemberRole = "Master"
	LabMemberRoleBachelor   LabMemberRole = "Bachelor"
	LabMemberRoleResearcher LabMemberRole = "Researcher"
)

// ProjectStatus defines the possible statuses for projects
type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusCompleted ProjectStatus = "completed"
)
