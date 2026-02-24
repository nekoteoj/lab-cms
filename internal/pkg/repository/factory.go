// Package repository provides data access layer implementations for all entities.
// This file contains the factory for creating and initializing all repositories.
package repository

import (
	"github.com/nekoteoj/lab-cms/internal/pkg/db"
)

// Factory manages all repository instances and provides centralized access.
type Factory struct {
	DBManager        *db.DBManager
	Users            *UserRepository
	LabMembers       *LabMemberRepository
	Publications     *PublicationRepository
	Projects         *ProjectRepository
	News             *NewsRepository
	HomepageSections *HomepageRepository
}

// NewFactory creates and initializes all repositories with a shared database connection.
func NewFactory(dbManager *db.DBManager) *Factory {
	return &Factory{
		DBManager:        dbManager,
		Users:            NewUserRepository(dbManager),
		LabMembers:       NewLabMemberRepository(dbManager),
		Publications:     NewPublicationRepository(dbManager),
		Projects:         NewProjectRepository(dbManager),
		News:             NewNewsRepository(dbManager),
		HomepageSections: NewHomepageRepository(dbManager),
	}
}

// Close closes the database connection.
// Should be called during graceful shutdown.
func (f *Factory) Close() error {
	return f.DBManager.Close()
}
