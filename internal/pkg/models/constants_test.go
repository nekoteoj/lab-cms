package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants_UserRole(t *testing.T) {
	tests := []struct {
		role     UserRole
		expected string
	}{
		{UserRoleNormal, "normal"},
		{UserRoleRoot, "root"},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.role))
		})
	}
}

func TestConstants_LabMemberRole(t *testing.T) {
	tests := []struct {
		role     LabMemberRole
		expected string
	}{
		{LabMemberRolePI, "PI"},
		{LabMemberRolePostdoc, "Postdoc"},
		{LabMemberRolePhD, "PhD"},
		{LabMemberRoleMaster, "Master"},
		{LabMemberRoleBachelor, "Bachelor"},
		{LabMemberRoleResearcher, "Researcher"},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.role))
		})
	}
}

func TestConstants_ProjectStatus(t *testing.T) {
	tests := []struct {
		status   ProjectStatus
		expected string
	}{
		{ProjectStatusActive, "active"},
		{ProjectStatusCompleted, "completed"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestConstants_HomepageSectionKeys(t *testing.T) {
	tests := []struct {
		key      string
		expected string
	}{
		{HomepageSectionOverview, "overview"},
		{HomepageSectionMission, "mission"},
		{HomepageSectionResearch, "research"},
		{HomepageSectionContact, "contact"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.key)
		})
	}
}
