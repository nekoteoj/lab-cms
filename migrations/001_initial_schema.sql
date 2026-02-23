-- Initial database schema for Lab CMS
-- Creates all core tables including users, lab members, publications, projects, news, and homepage content
-- Also creates junction tables for many-to-many relationships

-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- Users table: stores admin authentication information
-- Roles: 'normal' for regular admins, 'root' for root administrators
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('normal', 'root')) DEFAULT 'normal',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Lab members table: stores information about lab personnel
-- Members are organized by role and can be marked as alumni
CREATE TABLE lab_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('PI', 'Postdoc', 'PhD', 'Master', 'Bachelor', 'Researcher')),
    email TEXT,
    bio TEXT,
    photo_url TEXT,
    personal_page_content TEXT,
    research_interests TEXT,
    is_alumni BOOLEAN DEFAULT 0,
    display_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Publications table: stores lab research publications
-- Authors stored as text for flexibility, with junction table for linking to members
CREATE TABLE publications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    authors_text TEXT NOT NULL,
    venue TEXT,
    year INTEGER NOT NULL,
    url TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Projects table: stores research projects
-- Status indicates if project is active or completed
CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('active', 'completed')) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- News table: stores lab announcements and news items
-- Supports draft/published workflow with published_at timestamp
CREATE TABLE news (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    published_at DATETIME,
    is_published BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Homepage sections table: stores editable content sections for the homepage
-- Each section identified by a unique key (e.g., 'overview', 'mission', etc.)
CREATE TABLE homepage_sections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    section_key TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    display_order INTEGER DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Junction table: links projects to lab members
-- Represents which members are part of which projects
CREATE TABLE project_members (
    project_id INTEGER NOT NULL,
    member_id INTEGER NOT NULL,
    PRIMARY KEY (project_id, member_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES lab_members(id) ON DELETE CASCADE
);

-- Junction table: links publications to lab members
-- Represents which members are authors of which publications
CREATE TABLE publication_authors (
    publication_id INTEGER NOT NULL,
    member_id INTEGER NOT NULL,
    PRIMARY KEY (publication_id, member_id),
    FOREIGN KEY (publication_id) REFERENCES publications(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES lab_members(id) ON DELETE CASCADE
);

-- Junction table: links projects to publications
-- Represents which publications are associated with which projects
CREATE TABLE project_publications (
    project_id INTEGER NOT NULL,
    publication_id INTEGER NOT NULL,
    PRIMARY KEY (project_id, publication_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (publication_id) REFERENCES publications(id) ON DELETE CASCADE
);
