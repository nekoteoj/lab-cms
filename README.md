# Lab CMS

A content management system for research labs implemented in Go.

## Features

- Full-stack multi-page web application
- SQLite database for portability
- HTML templating
- Environment-based configuration

## Project Structure

```
lab-cms/
├── cmd/                    # Main applications
│   └── server/             # Web server entry point
├── internal/               # Private application code
│   ├── app/                # Application-specific code
│   └── pkg/                # Shared internal packages
│       ├── config/         # Configuration loading
│       ├── models/         # Database models
│       ├── repository/     # Data access layer
│       └── services/       # Business logic
├── web/                    # Web assets
│   ├── static/             # CSS, JS, images
│   └── templates/          # HTML templates
├── migrations/             # Database migrations
├── configs/                # Configuration templates
├── test/                   # Test data and helpers
├── scripts/                # Build and utility scripts
├── build/                  # Packaging and CI
├── deployments/            # Deployment configurations
└── docs/                   # Documentation
```

## Quick Start

```bash
# Copy environment configuration
cp configs/.env.example .env

# Run the server
make run
```

## Requirements

- Go 1.21 or later

## License

See [LICENSE](LICENSE) for details.
