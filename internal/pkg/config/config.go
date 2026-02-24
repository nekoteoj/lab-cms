// Package config provides configuration management for the Lab CMS application.
// Configuration is loaded from environment variables with sensible defaults.
//
// Security defaults are prioritized:
//   - CSRF protection enabled by default
//   - Session cookies use HttpOnly and SameSite=Strict
//   - Secure cookies in production mode
//   - Session secret required (no default for security)
//   - Production mode enforces stricter security rules
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Server configuration
	Port string // Server port (default: 8080)
	Env  string // Environment: development, production (default: development)

	// Database configuration
	DatabaseURL    string // SQLite database file path (default: ./data/lab-cms.db)
	DBMaxOpenConns int    // Maximum number of open connections (default: 0 = unlimited)
	DBMaxIdleConns int    // Maximum number of idle connections (default: 0 = Go default)

	// Session & Security
	SessionSecret  string // Required: Secret for session signing (no default)
	SessionMaxAge  int    // Session lifetime in hours (default: 24)
	CookieSecure   bool   // HTTPS only cookies (default: false in dev, true in prod)
	CookieHttpOnly bool   // Prevent JavaScript access to cookies (default: true)
	CookieSameSite string // CSRF protection: strict, lax, none (default: strict)
	CSRFEnabled    bool   // Enable CSRF token validation (default: true)
	TrustedProxies string // Comma-separated list of trusted proxy IPs (default: empty)

	// Initial admin setup (one-time use for first deployment)
	RootAdminUsername string // Username for initial root admin (default: admin)
	RootAdminPassword string // Password for initial root admin (default: empty - must be set)

	// Upload configuration
	UploadPath    string // Directory for file uploads (default: ./uploads)
	MaxUploadSize int64  // Maximum file upload size in bytes (default: 10485760 = 10MB)

	// Logging
	LogLevel string // Log level: debug, info, warn, error (default: info)
}

// Load reads configuration from environment variables and .env file.
// It applies sensible defaults and returns a Config struct.
// Note: Call Validate() to ensure required values are set.
func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:              getEnv("PORT", "8080"),
		Env:               getEnv("ENV", "development"),
		DatabaseURL:       getEnv("DATABASE_URL", "./data/lab-cms.db"),
		DBMaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 0), // 0 = use Go default (unlimited)
		DBMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 0), // 0 = use Go default (2)
		SessionSecret:     getEnv("SESSION_SECRET", ""),
		SessionMaxAge:     getEnvInt("SESSION_MAX_AGE", 24),
		CookieSecure:      getEnvBool("COOKIE_SECURE", false),
		CookieHttpOnly:    getEnvBool("COOKIE_HTTPONLY", true),
		CookieSameSite:    getEnv("COOKIE_SAMESITE", "strict"),
		CSRFEnabled:       getEnvBool("CSRF_ENABLED", true),
		TrustedProxies:    getEnv("TRUSTED_PROXIES", ""),
		RootAdminUsername: getEnv("ROOT_ADMIN_USERNAME", "admin"),
		RootAdminPassword: getEnv("ROOT_ADMIN_PASSWORD", ""),
		UploadPath:        getEnv("UPLOAD_PATH", "./uploads"),
		MaxUploadSize:     getEnvInt64("MAX_UPLOAD_SIZE", 10485760), // 10MB
		LogLevel:          strings.ToLower(getEnv("LOG_LEVEL", "info")),
	}

	// Auto-enable secure cookies in production
	if cfg.Env == "production" {
		cfg.CookieSecure = true
	}

	return cfg
}

// Validate checks that all required configuration values are set correctly.
// Returns an error if validation fails with details about what needs to be fixed.
func (c *Config) Validate() error {
	var errors []string

	// Validate required fields
	if c.SessionSecret == "" {
		errors = append(errors, "SESSION_SECRET is required - generate with: openssl rand -base64 32")
	}

	if c.RootAdminPassword == "" {
		errors = append(errors, "ROOT_ADMIN_PASSWORD is required for initial setup (minimum 8 characters)")
	} else if len(c.RootAdminPassword) < 8 {
		errors = append(errors, "ROOT_ADMIN_PASSWORD must be at least 8 characters")
	}

	// Validate port is numeric
	if _, err := strconv.Atoi(c.Port); err != nil {
		errors = append(errors, fmt.Sprintf("PORT must be a valid number, got: %s", c.Port))
	}

	// Validate environment value
	if c.Env != "development" && c.Env != "production" {
		errors = append(errors, fmt.Sprintf("ENV must be 'development' or 'production', got: %s", c.Env))
	}

	// Validate log level
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.LogLevel] {
		errors = append(errors, fmt.Sprintf("LOG_LEVEL must be debug, info, warn, or error, got: %s", c.LogLevel))
	}

	// Validate session max age is positive
	if c.SessionMaxAge <= 0 {
		errors = append(errors, "SESSION_MAX_AGE must be a positive number of hours")
	}

	// Validate SameSite value
	validSameSite := map[string]bool{"strict": true, "lax": true, "none": true}
	if !validSameSite[strings.ToLower(c.CookieSameSite)] {
		errors = append(errors, fmt.Sprintf("COOKIE_SAMESITE must be strict, lax, or none, got: %s", c.CookieSameSite))
	}

	// Validate upload path exists or can be created
	if c.UploadPath != "" {
		if err := ensureDir(c.UploadPath); err != nil {
			errors = append(errors, fmt.Sprintf("UPLOAD_PATH directory cannot be created: %v", err))
		}
	}

	// Production-specific security checks
	if c.Env == "production" {
		if len(c.SessionSecret) < 32 {
			errors = append(errors, "SESSION_SECRET must be at least 32 characters in production")
		}

		if c.CSRFEnabled == false {
			errors = append(errors, "CSRF_ENABLED cannot be false in production")
		}

		if c.CookieHttpOnly == false {
			errors = append(errors, "COOKIE_HTTPONLY cannot be false in production")
		}

		if strings.ToLower(c.CookieSameSite) != "strict" {
			errors = append(errors, "COOKIE_SAMESITE must be 'strict' in production")
		}

		if c.LogLevel == "debug" {
			errors = append(errors, "LOG_LEVEL cannot be 'debug' in production")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed:\n- %s", strings.Join(errors, "\n- "))
	}

	return nil
}

// IsProduction returns true if the application is running in production mode.
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// IsDevelopment returns true if the application is running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		lower := strings.ToLower(value)
		if lower == "true" || lower == "1" || lower == "yes" || lower == "on" {
			return true
		}
		if lower == "false" || lower == "0" || lower == "no" || lower == "off" {
			return false
		}
	}
	return defaultValue
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0750)
}
