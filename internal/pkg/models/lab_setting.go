package models

import (
	"time"
)

// LabSetting represents a key-value pair for lab configuration
// Used to store configurable settings like lab name and description
type LabSetting struct {
	ID           int       `json:"id"`
	SettingKey   string    `json:"setting_key" validate:"required,max=255"`
	SettingValue string    `json:"setting_value" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Common lab setting keys
const (
	LabSettingName        = "lab_name"
	LabSettingDescription = "lab_description"
)
