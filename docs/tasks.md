# Project Tasks - Lab CMS

Task board for Lab CMS development. Organized by milestones with Trello-style columns.

---

## Milestones

- **MVP v0.1** - Core Infrastructure & Foundation
- **MVP v0.2** - Admin System & Content Management
- **MVP v0.3** - Public Website
- **Beta v0.9** - Testing, Security & Polish
- **v1.0** - Production Ready

---

## Board Columns

### Backlog
Tasks planned for future development, not yet prioritized.

### To Do
Ready for development, prioritized by milestone.

### In Progress
Currently being worked on by an agent.

### Review
Code complete, needs code review and/or testing.

### Done
Completed, reviewed, and tested.

---

## Tasks

---

### [INF-001] Project Setup & Dependencies
**Status:** Done  
**Milestone:** MVP v0.1  
**Priority:** High  
**Story Points:** 2  
**Assignee:** opencode  
**Dependencies:** None

**Description:**  
Initialize Go project structure, set up module, install required dependencies (sqlite driver, godotenv, etc.), create Makefile with basic commands.

**Acceptance Criteria:**
- [x] `go.mod` initialized with project name
- [x] All dependencies added and `go mod tidy` passes
- [x] Makefile created with run, build, test, clean commands
- [x] Project directory structure matches AGENTS.md specification
- [x] `.gitignore` configured for Go projects

---

### [INF-002] Database Schema Design & Migrations
**Status:** Done  
**Milestone:** MVP v0.1  
**Priority:** High  
**Story Points:** 3  
**Assignee:** opencode  
**Dependencies:** INF-001

**Description:**  
Design database schema for users, lab members, publications, projects, news, and homepage content. Create migration files.

**Acceptance Criteria:**
- [x] Schema designed for all required entities
- [x] Migration files created in `migrations/` directory
- [x] Migration runner implemented
- [x] All tables have proper indexes
- [x] Foreign key relationships defined
- [x] Schema documented in code comments

---

### [INF-003] Base Models
**Status:** Done  
**Milestone:** MVP v0.1  
**Priority:** High  
**Story Points:** 2  
**Assignee:** opencode  
**Dependencies:** INF-002

**Description:**  
Create Go structs for all database entities in `internal/pkg/models/`. Include proper JSON tags and validation tags.

**Acceptance Criteria:**
- [x] User model (id, email, role, created_at, updated_at) - password handled separately
- [x] LabMember model (id, name, role, email, bio, photo_url, personal_page_content, research_interests, is_alumni, display_order, created_at, updated_at)
- [x] Publication model (id, title, authors_text, venue, year, url, created_at, updated_at)
- [x] Project model (id, title, description, status, created_at, updated_at)
- [x] News model (id, title, content, published_at, is_published, created_at, updated_at)
- [x] HomepageSection model (id, section_key, title, content, display_order, updated_at)
- [x] All models have proper field tags (JSON and validation)
- [x] Constants defined for enums (UserRole, LabMemberRole, ProjectStatus)

---

### [INF-004] Configuration Management
**Status:** Done  
**Milestone:** MVP v0.1  
**Priority:** High  
**Story Points:** 2  
**Assignee:** opencode  
**Dependencies:** INF-001

**Description:**  
Implement configuration loading from environment variables with sensible defaults. Create `.env.example` template.

**Acceptance Criteria:**
- [x] Config struct defined with all required fields
- [x] Environment variables loaded with godotenv
- [x] Sensible defaults for development
- [x] Validation of required config values
- [x] `.env.example` created in `configs/`
- [x] Configuration documented

---

### [INF-005] Base Repository Layer
**Status:** Done  
**Milestone:** MVP v0.1  
**Priority:** High  
**Story Points:** 3  
**Assignee:** opencode  
**Dependencies:** INF-003, INF-004

**Description:**  
Create repository pattern with database connection handling. Implement base repository with common CRUD operations.

**Acceptance Criteria:**
- [x] Database connection pool setup
- [x] Base repository interface defined
- [x] Transaction support implemented
- [x] Connection error handling
- [x] Repository initialization in app startup
- [x] All SQL uses parameterized queries

---

### [INF-008] Lab Settings Database Schema
**Status:** To Do  
**Milestone:** MVP v0.1  
**Priority:** High  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** INF-002

**Description:**  
Create database migration for lab_settings table with key-value structure. Insert default values during migration.

**Acceptance Criteria:**
- [ ] Migration file for `lab_settings` table (id, setting_key, setting_value, created_at, updated_at)
- [ ] Unique constraint on setting_key
- [ ] Default values inserted: lab_name="Research Lab", lab_description="A research laboratory"
- [ ] Migration runs successfully with `make test`
- [ ] LabSetting model created in `internal/pkg/models/`

---

### [INF-006] Error Handling & Logging Framework
**Status:** To Do  
**Milestone:** MVP v0.1  
**Priority:** Medium  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** INF-001

**Description:**  
Set up consistent error handling patterns and logging framework across the application.

**Acceptance Criteria:**
- [ ] Custom error types defined
- [ ] Error wrapping with context
- [ ] Structured logging setup
- [ ] Log levels configured (debug, info, error)
- [ ] HTTP error response helpers
- [ ] Error handling middleware

---

### [INF-007] Base HTTP Server Setup
**Status:** To Do  
**Milestone:** MVP v0.1  
**Priority:** High  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** INF-004, INF-006

**Description:**  
Set up HTTP server with routing, middleware chain, and graceful shutdown handling.

**Acceptance Criteria:**
- [ ] HTTP server initialized with configurable port
- [ ] Router setup (gorilla/mux or standard library)
- [ ] Middleware chain (logging, recovery, security headers)
- [ ] Graceful shutdown on SIGTERM/SIGINT
- [ ] Health check endpoint
- [ ] Server startup/shutdown logging

---

### [AUTH-001] Authentication System
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** High  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** INF-005, INF-007

**Description:**  
Implement user login/logout functionality with session-based authentication.

**Acceptance Criteria:**
- [ ] Login form handler (GET/POST)
- [ ] Logout handler
- [ ] Password verification against hash
- [ ] Session creation on successful login
- [ ] Login error messages (don't leak if user exists)
- [ ] Login page HTML template
- [ ] Redirect to admin dashboard on success

---

### [AUTH-002] Session Management
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** High  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** AUTH-001

**Description:**  
Implement secure session management with session store (database or cookie-based).

**Acceptance Criteria:**
- [ ] Session store implementation
- [ ] Session cookie configuration (HttpOnly, Secure, SameSite)
- [ ] Session expiration handling
- [ ] Session retrieval middleware
- [ ] Session destruction on logout
- [ ] CSRF token generation per session

---

### [AUTH-003] Password Security
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** High  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** INF-003

**Description:**  
Implement secure password hashing using bcrypt. Create password validation utilities.

**Acceptance Criteria:**
- [ ] Password hashing with bcrypt (cost factor 10+)
- [ ] Password verification function
- [ ] Password strength validation (optional but recommended)
- [ ] Secure password comparison (timing-safe)
- [ ] Unit tests for hash/verify functions

---

### [USER-001] User Role System
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** High  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** INF-003, AUTH-002

**Description:**  
Implement role-based access control with Normal Admin and Root Admin roles.

**Acceptance Criteria:**
- [ ] Role constants defined (normal, root)
- [ ] Role-based middleware for authorization
- [ ] Helper functions to check roles
- [ ] Root admin can manage other admins
- [ ] Normal admin cannot access user management
- [ ] Role stored in session

---

### [ADMIN-001] Admin Dashboard Layout
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** High  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** INF-007, AUTH-002

**Description:**  
Create admin dashboard layout with navigation sidebar, header, and content area. Base HTML templates.

**Acceptance Criteria:**
- [ ] Admin base layout template
- [ ] Navigation sidebar with menu items
- [ ] Header with user info and logout
- [ ] Responsive layout (mobile-friendly)
- [ ] Active menu item highlighting
- [ ] Consistent styling across admin pages

---

### [ADMIN-009] Lab Settings Management
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** High  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** INF-008, ADMIN-001

**Description:**  
Implement lab settings management UI. Allow root admins to configure lab name and description.

**Acceptance Criteria:**
- [ ] Repository methods for get/update settings in `internal/pkg/repository/`
- [ ] Admin settings page with form (root admin only)
- [ ] Form fields: lab name, lab description
- [ ] Display current values in form
- [ ] Save changes with validation
- [ ] Flash messages for success/error
- [ ] Settings used by public site header and homepage

---

### [ADMIN-002] Homepage Content Management
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** High  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** ADMIN-001, INF-005

**Description:**  
Implement CRUD for homepage content sections (lab overview, featured content).

**Acceptance Criteria:**
- [ ] Homepage content list view
- [ ] Edit form for homepage sections
- [ ] WYSIWYG or markdown editor for content
- [ ] Preview functionality
- [ ] Save changes to database
- [ ] Success/error flash messages

---

### [ADMIN-003] Lab Member CRUD
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** High  
**Story Points:** 5  
**Assignee:** (unassigned)  
**Dependencies:** ADMIN-001, INF-005

**Description:**  
Implement full CRUD for lab members with role types (PI, Postdoc, PhD, Master, Bachelor, Researcher, Alumni).

**Acceptance Criteria:**
- [ ] Member list view with filtering by role
- [ ] Add new member form with all fields
- [ ] Photo upload with preview
- [ ] Edit member form
- [ ] Delete member (or move to alumni)
- [ ] Role selection dropdown
- [ ] Personal page content editor
- [ ] Research interests field
- [ ] Validation on all fields

---

### [ADMIN-004] Publication Management
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** Medium  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** ADMIN-001, INF-005

**Description:**  
Implement CRUD for lab publications with simple listing format.

**Acceptance Criteria:**
- [ ] Publication list view (chronological)
- [ ] Add new publication form
- [ ] Fields: title, authors, venue, year, URL (optional)
- [ ] Edit publication form
- [ ] Delete publication
- [ ] Link to related members/projects
- [ ] Validation on required fields

---

### [ADMIN-005] Research Project Management
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** Medium  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** ADMIN-001, INF-005, ADMIN-003

**Description:**  
Implement CRUD for research projects with team member assignment and publication linking.

**Acceptance Criteria:**
- [ ] Project list view with status
- [ ] Add new project form
- [ ] Fields: title, description, status (active/completed)
- [ ] Assign team members (multi-select)
- [ ] Link publications to project
- [ ] Edit project form
- [ ] Delete project
- [ ] Project status toggle

---

### [ADMIN-006] News Management
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** Medium  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** ADMIN-001, INF-005

**Description:**  
Implement CRUD for news and announcements with publication date.

**Acceptance Criteria:**
- [ ] News list view (chronological)
- [ ] Add news form
- [ ] Fields: title, content, publish date
- [ ] Edit news form
- [ ] Delete news
- [ ] Published/draft status
- [ ] Archive old news

---

### [ADMIN-007] Root Admin User Management
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** High  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** USER-001, ADMIN-001

**Description:**  
Implement user management for root admins only - add/remove admins, assign roles, reset passwords.

**Acceptance Criteria:**
- [ ] User list view (admin only)
- [ ] Add new admin form with role selection
- [ ] Edit admin (change role, reset password)
- [ ] Delete admin account
- [ ] Password reset functionality
- [ ] Root admin cannot delete themselves
- [ ] Protected by role middleware

---

### [ADMIN-008] File Upload System
**Status:** To Do  
**Milestone:** MVP v0.2  
**Priority:** Medium  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** ADMIN-001

**Description:**  
Implement secure file upload system for member photos and other assets.

**Acceptance Criteria:**
- [ ] File upload handler
- [ ] Image validation (type, size)
- [ ] Image resizing/thumbnails (optional)
- [ ] Secure storage location (outside web root)
- [ ] File serving endpoint
- [ ] File naming (UUID to prevent collisions)
- [ ] Cleanup on delete

---

### [PUB-001] Public Layout & Navigation
**Status:** To Do  
**Milestone:** MVP v0.3  
**Priority:** High  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** INF-007

**Description:**  
Create public website layout with header, navigation, footer, and content area.

**Acceptance Criteria:**
- [ ] Public base layout template
- [ ] Header with lab name/logo
- [ ] Navigation menu (Home, Members, Publications, Projects, News)
- [ ] Footer with copyright
- [ ] Responsive design
- [ ] Clean, professional styling

---

### [PUB-002] Homepage View
**Status:** To Do  
**Milestone:** MVP v0.3  
**Priority:** High  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** PUB-001, ADMIN-002

**Description:**  
Implement public homepage displaying lab overview and recent content.

**Acceptance Criteria:**
- [ ] Display lab overview from database
- [ ] Show recent news (3-5 items)
- [ ] Show featured projects
- [ ] Proper HTML rendering of content
- [ ] SEO meta tags
- [ ] Clean URL (/)

---

### [PUB-003] Lab Members Page
**Status:** To Do  
**Milestone:** MVP v0.3  
**Priority:** High  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** PUB-001, ADMIN-003

**Description:**  
Implement public members page showing all lab members organized by role with filtering.

**Acceptance Criteria:**
- [ ] Members grouped by role (PI, Postdoc, PhD, Master, Bachelor, Researcher)
- [ ] Each member shows photo, name, title, brief bio
- [ ] Link to member's personal page
- [ ] Alumni section (separate or at bottom)
- [ ] Responsive grid layout
- [ ] Clean URL (/members)

---

### [PUB-004] Member Profile Pages
**Status:** To Do  
**Milestone:** MVP v0.3  
**Priority:** Medium  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** PUB-003

**Description:**  
Implement individual member profile pages with detailed information.

**Acceptance Criteria:**
- [ ] Display full member info (photo, name, role, bio)
- [ ] Show research interests
- [ ] Show personal page content
- [ ] List member's publications
- [ ] List member's projects
- [ ] Clean URL (/members/:id or /members/:slug)

---

### [PUB-005] Publications Listing Page
**Status:** To Do  
**Milestone:** MVP v0.3  
**Priority:** Medium  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** PUB-001, ADMIN-004

**Description:**  
Implement public publications page showing all lab publications in chronological order.

**Acceptance Criteria:**
- [ ] List all publications (newest first)
- [ ] Show title, authors, venue, year
- [ ] Link to external URL if available
- [ ] Group by year (optional)
- [ ] Responsive layout
- [ ] Clean URL (/publications)

---

### [PUB-006] Research Projects Page
**Status:** To Do  
**Milestone:** MVP v0.3  
**Priority:** Medium  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** PUB-001, ADMIN-005

**Description:**  
Implement public projects page showing ongoing and completed research projects.

**Acceptance Criteria:**
- [ ] List projects with title, description, status
- [ ] Show team members on project
- [ ] Show linked publications
- [ ] Filter by status (active/completed)
- [ ] Responsive layout
- [ ] Clean URL (/projects)

---

### [PUB-007] News & Events Page
**Status:** To Do  
**Milestone:** MVP v0.3  
**Priority:** Low  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** PUB-001, ADMIN-006

**Description:**  
Implement public news page showing all lab news and announcements.

**Acceptance Criteria:**
- [ ] List news chronologically
- [ ] Show title, date, content preview
- [ ] Individual news item pages
- [ ] Pagination if needed
- [ ] Responsive layout
- [ ] Clean URL (/news)

---

### [PUB-008] Static Pages
**Status:** To Do  
**Milestone:** MVP v0.3  
**Priority:** Low  
**Story Points:** 1  
**Assignee:** (unassigned)  
**Dependencies:** PUB-001

**Description:**  
Create static pages: About, Contact, and 404 error page.

**Acceptance Criteria:**
- [ ] About page with lab description
- [ ] Contact page with lab contact info
- [ ] 404 error page with friendly message
- [ ] Clean URLs (/about, /contact)

---

### [TEST-001] Unit Tests - Models
**Status:** To Do  
**Milestone:** Beta v0.9  
**Priority:** High  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** INF-003

**Description:**  
Write comprehensive unit tests for all model structs and validation logic.

**Acceptance Criteria:**
- [ ] Tests for all model constructors
- [ ] Tests for validation functions
- [ ] Tests for model methods
- [ ] 80%+ code coverage for models
- [ ] All tests pass

---

### [TEST-002] Unit Tests - Repositories
**Status:** Done  
**Milestone:** Beta v0.9  
**Priority:** High  
**Story Points:** 5  
**Assignee:** opencode  
**Dependencies:** INF-005

**Description:**  
Write unit tests for repository layer using in-memory SQLite or mocks.

**Acceptance Criteria:**
- [x] Tests for all repository methods
- [x] Test CRUD operations for each entity
- [x] Test error scenarios
- [x] 80%+ code coverage for repositories
- [x] All tests pass

---

### [TEST-003] Unit Tests - Handlers
**Status:** To Do  
**Milestone:** Beta v0.9  
**Priority:** High  
**Story Points:** 5  
**Assignee:** (unassigned)  
**Dependencies:** All handler implementations

**Description:**  
Write unit tests for HTTP handlers using httptest and mocks.

**Acceptance Criteria:**
- [ ] Tests for all HTTP handlers
- [ ] Test request/response cycles
- [ ] Test authentication middleware
- [ ] Test error handling
- [ ] 70%+ code coverage for handlers
- [ ] All tests pass

---

### [TEST-004] Integration Tests
**Status:** To Do  
**Milestone:** Beta v0.9  
**Priority:** Medium  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** TEST-001, TEST-002, TEST-003

**Description:**  
Write integration tests for end-to-end workflows (login → create content → view public).

**Acceptance Criteria:**
- [ ] Full authentication flow test
- [ ] Content creation to public view test
- [ ] Database transaction tests
- [ ] File upload/download test
- [ ] All integration tests pass

---

### [SEC-001] Security Review
**Status:** To Do  
**Milestone:** Beta v0.9  
**Priority:** High  
**Story Points:** 3  
**Assignee:** (unassigned)  
**Dependencies:** All implementation tasks

**Description:**  
Perform security audit focusing on common vulnerabilities.

**Acceptance Criteria:**
- [ ] SQL Injection prevention verified (parameterized queries)
- [ ] XSS prevention verified (HTML escaping in templates)
- [ ] CSRF protection implemented and tested
- [ ] Authentication bypass attempts tested
- [ ] Authorization checks verified
- [ ] Session fixation/hijacking protection
- [ ] File upload security reviewed
- [ ] Security headers configured
- [ ] HTTPS redirect (if applicable)

---

### [UI-001] Responsive Design Review
**Status:** To Do  
**Milestone:** Beta v0.9  
**Priority:** Medium  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** All UI tasks

**Description:**  
Review and test responsive design across devices and screen sizes.

**Acceptance Criteria:**
- [ ] Test on mobile (320px+)
- [ ] Test on tablet (768px+)
- [ ] Test on desktop (1024px+)
- [ ] Admin dashboard usable on mobile
- [ ] Public site looks good on all sizes
- [ ] Images scale properly
- [ ] Navigation works on mobile

---

### [DOC-001] Deployment Documentation
**Status:** To Do  
**Milestone:** Beta v0.9  
**Priority:** Medium  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** All implementation tasks

**Description:**  
Create deployment guide with environment setup, database migration, and server configuration.

**Acceptance Criteria:**
- [ ] Environment setup instructions
- [ ] Database migration steps
- [ ] Build and run instructions
- [ ] Production configuration guide
- [ ] Troubleshooting section
- [ ] Security hardening checklist

---

### [DEPLOY-001] Production Build Setup
**Status:** To Do  
**Milestone:** v1.0  
**Priority:** High  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** All Beta tasks

**Description:**  
Set up production build process, optimization, and binary configuration.

**Acceptance Criteria:**
- [ ] Production build command
- [ ] Optimized binary size
- [ ] No debug symbols in production
- [ ] Asset embedding (if needed)
- [ ] Build script in Makefile
- [ ] Cross-compilation support (optional)

---

### [DEPLOY-002] Production Configuration
**Status:** To Do  
**Milestone:** v1.0  
**Priority:** High  
**Story Points:** 2  
**Assignee:** (unassigned)  
**Dependencies:** DEPLOY-001

**Description:**  
Create production environment configuration templates and scripts.

**Acceptance Criteria:**
- [ ] Production .env template
- [ ] Database backup script
- [ ] Log rotation configuration
- [ ] Systemd service file (optional)
- [ ] Docker configuration (optional)
- [ ] Health check endpoint working

---

### [REL-001] Final QA & Bug Fixes
**Status:** To Do  
**Milestone:** v1.0  
**Priority:** High  
**Story Points:** 5  
**Assignee:** (unassigned)  
**Dependencies:** All previous tasks

**Description:**  
Final quality assurance testing and bug fixes before production release.

**Acceptance Criteria:**
- [ ] All tests passing
- [ ] Manual testing completed
- [ ] All critical bugs fixed
- [ ] Performance acceptable
- [ ] Documentation complete
- [ ] Ready for production deployment

---

## How to Use This Board

1. **Before starting work:**
   - Check "To Do" column for ready tasks
   - Verify dependencies are in "Done" or "Review"
   - Confirm no one else is working on it

2. **When starting a task:**
   - Move task to "In Progress"
   - Add your identifier as Assignee
   - Create feature branch: `git checkout -b feature/TASK-ID`

3. **When task is complete:**
   - Move task to "Review"
   - Ensure all acceptance criteria are checked
   - Run tests: `make test`
   - Request code review

4. **When reviewed:**
   - Move task to "Done"
   - Merge feature branch to main

5. **Adding new tasks:**
   - Add to "Backlog" with next available ID
   - Fill in all fields
   - Prioritize during planning

---

## Task ID Prefixes

- **INF-** - Infrastructure and setup
- **AUTH-** - Authentication and authorization
- **USER-** - User management
- **ADMIN-** - Admin system features
- **PUB-** - Public website features
- **TEST-** - Testing tasks
- **SEC-** - Security tasks
- **UI-** - UI/UX improvements
- **DOC-** - Documentation
- **DEPLOY-** - Deployment and release
- **REL-** - Release tasks

---

## Estimation Guide (Story Points)

- **1 point** - Trivial task (< 2 hours)
- **2 points** - Small task (2-4 hours)
- **3 points** - Medium task (4-8 hours)
- **5 points** - Large task (1-2 days)
- **8 points** - Very large task (2-3 days) - Consider breaking down
