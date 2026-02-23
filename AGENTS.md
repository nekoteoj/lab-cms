# AGENTS.md - Coding Guidelines for Lab CMS

This document provides guidelines for AI agents working on the Lab CMS codebase.

## Documentation Standards

### System Requirements
Functional requirements describing application features and behaviors belong in `docs/requirements.md`.

This document should describe:
- User roles and permissions
- Features and functionality the application supports
- User workflows and interactions
- Business logic and rules

This document should NOT describe:
- Technical implementation details (programming languages, frameworks, databases)
- System architecture or deployment configurations
- Code patterns or development tools

When adding new features, update the requirements document to reflect the new functionality.

## Task Management

All development tasks are tracked in `docs/tasks.md` using a Trello-style board format.

### How to Use

**Before starting work:**
- Check "To Do" column for ready tasks
- Verify dependencies are in "Done" or "Review"
- Confirm no one else is working on it (check Assignee field)

**When starting a task:**
- Move task to "In Progress" in `docs/tasks.md`
- Add your identifier as Assignee
- Create feature branch: `git checkout -b feature/TASK-ID` (e.g., `feature/INF-001`)

**When task is complete:**
- Move task to "Review" in `docs/tasks.md`
- Ensure all acceptance criteria are checked
- Run tests: `make test`
- Run linting: `golangci-lint run`
- Request code review

**When reviewed:**
- Move task to "Done" in `docs/tasks.md`
- Merge feature branch to main

**Adding new tasks:**
- Add to "Backlog" with next available ID
- Follow task ID prefixes defined in `docs/tasks.md`
- Include milestone, priority, story points, and acceptance criteria

### Task ID Prefixes

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

### Estimation Guide (Story Points)

- **1 point** - Trivial task (< 2 hours)
- **2 points** - Small task (2-4 hours)
- **3 points** - Medium task (4-8 hours)
- **5 points** - Large task (1-2 days)
- **8 points** - Very large task (2-3 days) - Consider breaking down

## Project Overview

- **Language**: Go 1.25.6
- **Framework**: Standard library (net/http), minimal external dependencies
- **Database**: SQLite via modernc.org/sqlite (CGO-free, pure Go)
- **Configuration**: Environment variables with godotenv support

## Build Commands

```bash
# Run the development server
make run
# Equivalent: go run ./cmd/server

# Build production binary
make build
# Equivalent: go build -o bin/server ./cmd/server

# Run all tests
make test
# Equivalent: go test ./...

# Run specific test
go test -v -run TestFunctionName ./path/to/package

# Clean build artifacts
make clean

# Download/update dependencies
go mod tidy

# Lint code (requires golangci-lint)
golangci-lint run

# Format code
goimports -w .
```

## Code Style Guidelines

### Imports
- Group imports in three sections:
  1. Standard library
  2. External dependencies
  3. Internal packages
- Use `goimports` for automatic import management
- Never use dot imports

### Formatting
- Run `gofmt` or `goimports` before committing
- Use tabs for indentation (Go standard)
- Maximum line length: 100 characters
- No trailing whitespace

### Naming Conventions
- **Exported identifiers**: PascalCase (Config, Load, UserRepository)
- **Unexported identifiers**: camelCase (getEnv, databaseURL)
- **Interface names**: Noun with -er suffix (Reader, Writer, Repository)
- **Function receivers**: 1-2 letter abbreviation (c *Config, db *Database)
- **Package names**: Short, lowercase, no underscores (config, models)
- **Constants**: PascalCase if exported, camelCase otherwise

### Error Handling
- Always check errors explicitly, never use `_` to ignore unless intentional
- Return errors rather than panicking in library code
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Use `log.Fatalf()` only in `cmd/server/main.go` for fatal startup errors

### Type Declarations
- Use structs with field tags for configuration and models
- Prefer value receivers for small structs, pointer receivers for large
- Example:
  ```go
  type Config struct {
      Port        string
      DatabaseURL string
      Env         string
  }
  
  func Load() *Config { ... }
  ```

### Project Structure
```
cmd/server/           # Application entry point (main.go only)
internal/
├── app/server/       # HTTP handlers, middleware, routes (private)
└── pkg/
    ├── config/       # Configuration loading and management
    ├── models/       # Database models and structs
    ├── repository/   # Data access layer (SQLite operations)
    └── services/     # Business logic layer
web/
├── static/           # CSS, JavaScript, images
└── templates/        # HTML templates (layouts/, pages/)
migrations/           # Database migration files
test/                 # Test data and helpers
configs/              # Configuration templates (.env.example)
scripts/              # Build and utility scripts
```

### Database Patterns
- Use `modernc.org/sqlite` driver (pure Go, no CGO required)
- Repository pattern: all SQL in `internal/pkg/repository/`
- Models in `internal/pkg/models/` (plain structs, no DB logic)
- Always use prepared statements or query parameters

### Configuration
- All configuration via environment variables
- Use `internal/pkg/config` package for loading
- Provide sensible defaults for development
- Never commit `.env` files or real credentials
- Example in `configs/.env.example`:
  ```
  PORT=8080
  ENV=development
  DATABASE_URL=./data/lab-cms.db
  ```

### Testing
- Test files: `*_test.go` alongside source files
- Use table-driven tests for multiple test cases
- Mock external dependencies (database, HTTP clients)
- Run single test: `go test -v -run TestFunctionName ./package`
- Test coverage target: 70% minimum

### Linting
Use `golangci-lint` with these enabled linters:
- errcheck: Error return value checking
- gosimple: Simplification suggestions
- govet: Vet reports
- ineffassign: Ineffective assignments
- staticcheck: Comprehensive analysis
- unused: Unused code detection

### Git Workflow
- Create feature branches: `git checkout -b feature/description`
- Write descriptive commit messages (present tense, imperative mood)
- Format: `type: description` (feat:, fix:, refactor:, docs:, test:)
- Example: "feat: add user authentication middleware"
- Keep commits atomic and focused
- Rebase feature branches on main before merging

### Security
- Never hardcode secrets in source code
- Validate all user input
- Use parameterized queries (prepared statements)
- Implement proper CSRF protection for forms
- Sanitize HTML output to prevent XSS

### Comments
- Document all exported types and functions with GoDoc format
- Comments should explain "why", not "what" (code explains what)
- Example: `// Load reads configuration from environment variables with defaults.`

### Commit Message Format
```
type: short description (50 chars max)

Optional longer description explaining what and why,
not how. Wrap at 72 characters.

Types:
- feat: New feature
- fix: Bug fix
- refactor: Code restructuring without behavior change
- docs: Documentation changes
- test: Test additions or modifications
- chore: Maintenance tasks, dependencies
- style: Formatting changes (gofmt, imports)
```
