-- Add performance indexes for frequently queried fields
-- These indexes improve query performance for common operations

-- Users table indexes
-- Unique constraint on email for fast lookups and preventing duplicates
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- Lab members table indexes
-- Composite index for filtering by alumni status and sorting by display order
CREATE INDEX idx_lab_members_alumni_order ON lab_members(is_alumni, display_order);
-- Index for role-based filtering
CREATE INDEX idx_lab_members_role ON lab_members(role);

-- Publications table indexes
-- Composite index for chronological listing (year descending, then creation date)
CREATE INDEX idx_publications_year_created ON publications(year DESC, created_at DESC);

-- Projects table indexes
-- Index for status-based filtering (active vs completed)
CREATE INDEX idx_projects_status ON projects(status);

-- News table indexes
-- Composite index for published news listing (newest first)
CREATE INDEX idx_news_published_created ON news(is_published, published_at DESC, created_at DESC);

-- Homepage sections table indexes
-- Unique index on section_key for fast lookups
CREATE UNIQUE INDEX idx_homepage_section_key ON homepage_sections(section_key);
-- Index for sorting sections by display order
CREATE INDEX idx_homepage_display_order ON homepage_sections(display_order);
