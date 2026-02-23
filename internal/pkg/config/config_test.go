package config

import (
	"os"
	"strconv"
	"testing"
)

// TestLoad_DefaultValues verifies that Load() returns sensible defaults
func TestLoad_DefaultValues(t *testing.T) {
	// Clear all related environment variables
	clearEnvVars()

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}
	if cfg.Env != "development" {
		t.Errorf("Expected Env to be 'development', got '%s'", cfg.Env)
	}
	if cfg.DatabaseURL != "./data/lab-cms.db" {
		t.Errorf("Expected DatabaseURL to be './data/lab-cms.db', got '%s'", cfg.DatabaseURL)
	}
	if cfg.SessionMaxAge != 24 {
		t.Errorf("Expected SessionMaxAge to be 24, got %d", cfg.SessionMaxAge)
	}
	if cfg.CookieSecure != false {
		t.Errorf("Expected CookieSecure to be false in dev, got %v", cfg.CookieSecure)
	}
	if cfg.CookieHttpOnly != true {
		t.Errorf("Expected CookieHttpOnly to be true, got %v", cfg.CookieHttpOnly)
	}
	if cfg.CookieSameSite != "strict" {
		t.Errorf("Expected CookieSameSite to be 'strict', got '%s'", cfg.CookieSameSite)
	}
	if cfg.CSRFEnabled != true {
		t.Errorf("Expected CSRFEnabled to be true, got %v", cfg.CSRFEnabled)
	}
	if cfg.RootAdminUsername != "admin" {
		t.Errorf("Expected RootAdminUsername to be 'admin', got '%s'", cfg.RootAdminUsername)
	}
	if cfg.UploadPath != "./uploads" {
		t.Errorf("Expected UploadPath to be './uploads', got '%s'", cfg.UploadPath)
	}
	if cfg.MaxUploadSize != 10485760 {
		t.Errorf("Expected MaxUploadSize to be 10485760, got %d", cfg.MaxUploadSize)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Expected LogLevel to be 'info', got '%s'", cfg.LogLevel)
	}
}

// TestLoad_EnvironmentValues verifies that Load() reads from environment variables
func TestLoad_EnvironmentValues(t *testing.T) {
	clearEnvVars()

	os.Setenv("PORT", "3000")
	os.Setenv("ENV", "production")
	os.Setenv("DATABASE_URL", "/custom/path/db.db")
	os.Setenv("SESSION_SECRET", "test-secret-32-chars-long-ok")
	os.Setenv("SESSION_MAX_AGE", "48")
	os.Setenv("COOKIE_SECURE", "true")
	os.Setenv("COOKIE_HTTPONLY", "false")
	os.Setenv("COOKIE_SAMESITE", "lax")
	os.Setenv("CSRF_ENABLED", "false")
	os.Setenv("TRUSTED_PROXIES", "127.0.0.1,10.0.0.1")
	os.Setenv("ROOT_ADMIN_USERNAME", "root")
	os.Setenv("ROOT_ADMIN_PASSWORD", "testpass123")
	os.Setenv("UPLOAD_PATH", "/custom/uploads")
	os.Setenv("MAX_UPLOAD_SIZE", "20971520")
	os.Setenv("LOG_LEVEL", "debug")

	cfg := Load()

	if cfg.Port != "3000" {
		t.Errorf("Expected Port to be '3000', got '%s'", cfg.Port)
	}
	if cfg.Env != "production" {
		t.Errorf("Expected Env to be 'production', got '%s'", cfg.Env)
	}
	if cfg.DatabaseURL != "/custom/path/db.db" {
		t.Errorf("Expected DatabaseURL to be '/custom/path/db.db', got '%s'", cfg.DatabaseURL)
	}
	if cfg.SessionSecret != "test-secret-32-chars-long-ok" {
		t.Errorf("Expected SessionSecret to be set correctly")
	}
	if cfg.SessionMaxAge != 48 {
		t.Errorf("Expected SessionMaxAge to be 48, got %d", cfg.SessionMaxAge)
	}
	if cfg.CookieSecure != true {
		t.Errorf("Expected CookieSecure to be true, got %v", cfg.CookieSecure)
	}
	if cfg.CookieHttpOnly != false {
		t.Errorf("Expected CookieHttpOnly to be false, got %v", cfg.CookieHttpOnly)
	}
	if cfg.CookieSameSite != "lax" {
		t.Errorf("Expected CookieSameSite to be 'lax', got '%s'", cfg.CookieSameSite)
	}
	if cfg.CSRFEnabled != false {
		t.Errorf("Expected CSRFEnabled to be false, got %v", cfg.CSRFEnabled)
	}
	if cfg.TrustedProxies != "127.0.0.1,10.0.0.1" {
		t.Errorf("Expected TrustedProxies to be '127.0.0.1,10.0.0.1', got '%s'", cfg.TrustedProxies)
	}
	if cfg.RootAdminUsername != "root" {
		t.Errorf("Expected RootAdminUsername to be 'root', got '%s'", cfg.RootAdminUsername)
	}
	if cfg.RootAdminPassword != "testpass123" {
		t.Errorf("Expected RootAdminPassword to be set correctly")
	}
	if cfg.UploadPath != "/custom/uploads" {
		t.Errorf("Expected UploadPath to be '/custom/uploads', got '%s'", cfg.UploadPath)
	}
	if cfg.MaxUploadSize != 20971520 {
		t.Errorf("Expected MaxUploadSize to be 20971520, got %d", cfg.MaxUploadSize)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel to be 'debug', got '%s'", cfg.LogLevel)
	}
}

// TestLoad_ProductionCookieSecure verifies that production mode auto-enables secure cookies
func TestLoad_ProductionCookieSecure(t *testing.T) {
	clearEnvVars()
	os.Setenv("ENV", "production")
	os.Setenv("COOKIE_SECURE", "false") // User tries to disable, but production overrides

	cfg := Load()

	if !cfg.CookieSecure {
		t.Error("Expected CookieSecure to be true in production mode")
	}
}

// TestLoad_BoolVariations tests various boolean string formats
func TestLoad_BoolVariations(t *testing.T) {
	truthyValues := []string{"true", "TRUE", "True", "1", "yes", "YES", "on", "ON"}
	falsyValues := []string{"false", "FALSE", "False", "0", "no", "NO", "off", "OFF"}

	for _, val := range truthyValues {
		t.Run("truthy_"+val, func(t *testing.T) {
			clearEnvVars()
			os.Setenv("CSRF_ENABLED", val)
			cfg := Load()
			if !cfg.CSRFEnabled {
				t.Errorf("Expected CSRF_ENABLED='%s' to be true", val)
			}
		})
	}

	for _, val := range falsyValues {
		t.Run("falsy_"+val, func(t *testing.T) {
			clearEnvVars()
			os.Setenv("CSRF_ENABLED", val)
			cfg := Load()
			if cfg.CSRFEnabled {
				t.Errorf("Expected CSRF_ENABLED='%s' to be false", val)
			}
		})
	}
}

// TestConfig_Validate_Success verifies valid configuration passes
func TestConfig_Validate_Success(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "development",
		DatabaseURL:       "./data/lab-cms.db",
		SessionSecret:     "valid-secret-32-chars-minimum-req",
		SessionMaxAge:     24,
		CookieSecure:      false,
		CookieHttpOnly:    true,
		CookieSameSite:    "strict",
		CSRFEnabled:       true,
		RootAdminUsername: "admin",
		RootAdminPassword: "validpass8",
		UploadPath:        "./uploads",
		MaxUploadSize:     10485760,
		LogLevel:          "info",
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}
}

// TestConfig_Validate_MissingSessionSecret verifies missing session secret fails
func TestConfig_Validate_MissingSessionSecret(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "development",
		SessionSecret:     "",
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "strict",
		SessionMaxAge:     24,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with missing session secret")
	}
	if err != nil && !contains(err.Error(), "SESSION_SECRET") {
		t.Errorf("Expected error to mention SESSION_SECRET, got: %v", err)
	}
}

// TestConfig_Validate_MissingRootAdminPassword verifies missing admin password fails
func TestConfig_Validate_MissingRootAdminPassword(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "development",
		SessionSecret:     "valid-secret-32-chars-minimum-req",
		RootAdminPassword: "",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "strict",
		SessionMaxAge:     24,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with missing root admin password")
	}
	if err != nil && !contains(err.Error(), "ROOT_ADMIN_PASSWORD") {
		t.Errorf("Expected error to mention ROOT_ADMIN_PASSWORD, got: %v", err)
	}
}

// TestConfig_Validate_ShortPassword verifies password minimum length
func TestConfig_Validate_ShortPassword(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "development",
		SessionSecret:     "valid-secret-32-chars-minimum-req",
		RootAdminPassword: "short7", // Only 7 chars, needs 8
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "strict",
		SessionMaxAge:     24,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with short password")
	}
	if err != nil && !contains(err.Error(), "8 characters") {
		t.Errorf("Expected error to mention 8 character requirement, got: %v", err)
	}
}

// TestConfig_Validate_InvalidPort verifies invalid port number fails
func TestConfig_Validate_InvalidPort(t *testing.T) {
	cfg := &Config{
		Port:              "invalid",
		Env:               "development",
		SessionSecret:     "valid-secret-32-chars-minimum-req",
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "strict",
		SessionMaxAge:     24,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with invalid port")
	}
	if err != nil && !contains(err.Error(), "PORT") {
		t.Errorf("Expected error to mention PORT, got: %v", err)
	}
}

// TestConfig_Validate_InvalidEnv verifies invalid environment value fails
func TestConfig_Validate_InvalidEnv(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "invalid",
		SessionSecret:     "valid-secret-32-chars-minimum-req",
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "strict",
		SessionMaxAge:     24,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with invalid environment")
	}
	if err != nil && !contains(err.Error(), "ENV") {
		t.Errorf("Expected error to mention ENV, got: %v", err)
	}
}

// TestConfig_Validate_InvalidLogLevel verifies invalid log level fails
func TestConfig_Validate_InvalidLogLevel(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "development",
		SessionSecret:     "valid-secret-32-chars-minimum-req",
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "strict",
		SessionMaxAge:     24,
		LogLevel:          "invalid",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with invalid log level")
	}
	if err != nil && !contains(err.Error(), "LOG_LEVEL") {
		t.Errorf("Expected error to mention LOG_LEVEL, got: %v", err)
	}
}

// TestConfig_Validate_InvalidSameSite verifies invalid SameSite value fails
func TestConfig_Validate_InvalidSameSite(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "development",
		SessionSecret:     "valid-secret-32-chars-minimum-req",
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "invalid",
		SessionMaxAge:     24,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with invalid SameSite")
	}
	if err != nil && !contains(err.Error(), "SAMESITE") {
		t.Errorf("Expected error to mention SAMESITE, got: %v", err)
	}
}

// TestConfig_Validate_InvalidSessionMaxAge verifies invalid session max age fails
func TestConfig_Validate_InvalidSessionMaxAge(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "development",
		SessionSecret:     "valid-secret-32-chars-minimum-req",
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "strict",
		SessionMaxAge:     0,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with invalid session max age")
	}
	if err != nil && !contains(err.Error(), "MAX_AGE") {
		t.Errorf("Expected error to mention MAX_AGE, got: %v", err)
	}
}

// TestConfig_Validate_Production_ShortSecret verifies production requires long secret
func TestConfig_Validate_Production_ShortSecret(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "production",
		SessionSecret:     "short-secret", // Less than 32 chars
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "strict",
		SessionMaxAge:     24,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with short secret in production")
	}
	if err != nil && !contains(err.Error(), "32 characters") {
		t.Errorf("Expected error to mention 32 character requirement, got: %v", err)
	}
}

// TestConfig_Validate_Production_CSRFDisabled verifies production requires CSRF
func TestConfig_Validate_Production_CSRFDisabled(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "production",
		SessionSecret:     "this-is-a-valid-32-char-secret",
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    true,
		CSRFEnabled:       false,
		CookieSameSite:    "strict",
		SessionMaxAge:     24,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with CSRF disabled in production")
	}
	if err != nil && !contains(err.Error(), "CSRF_ENABLED") {
		t.Errorf("Expected error to mention CSRF_ENABLED, got: %v", err)
	}
}

// TestConfig_Validate_Production_HttpOnlyDisabled verifies production requires HttpOnly
func TestConfig_Validate_Production_HttpOnlyDisabled(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "production",
		SessionSecret:     "this-is-a-valid-32-char-secret",
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    false,
		CSRFEnabled:       true,
		CookieSameSite:    "strict",
		SessionMaxAge:     24,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with HttpOnly disabled in production")
	}
	if err != nil && !contains(err.Error(), "HTTPONLY") {
		t.Errorf("Expected error to mention HTTPONLY, got: %v", err)
	}
}

// TestConfig_Validate_Production_LaxSameSite verifies production requires strict SameSite
func TestConfig_Validate_Production_LaxSameSite(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "production",
		SessionSecret:     "this-is-a-valid-32-char-secret",
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "lax",
		SessionMaxAge:     24,
		LogLevel:          "info",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with lax SameSite in production")
	}
	if err != nil && !contains(err.Error(), "SAMESITE") {
		t.Errorf("Expected error to mention SAMESITE, got: %v", err)
	}
}

// TestConfig_Validate_Production_DebugLogLevel verifies production forbids debug logging
func TestConfig_Validate_Production_DebugLogLevel(t *testing.T) {
	cfg := &Config{
		Port:              "8080",
		Env:               "production",
		SessionSecret:     "this-is-a-valid-32-char-secret",
		RootAdminPassword: "validpass8",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "strict",
		SessionMaxAge:     24,
		LogLevel:          "debug",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation to fail with debug logging in production")
	}
	if err != nil && !contains(err.Error(), "LOG_LEVEL") {
		t.Errorf("Expected error to mention LOG_LEVEL, got: %v", err)
	}
}

// TestConfig_IsProduction verifies production detection
func TestConfig_IsProduction(t *testing.T) {
	prodCfg := &Config{Env: "production"}
	if !prodCfg.IsProduction() {
		t.Error("Expected IsProduction() to return true for production env")
	}

	devCfg := &Config{Env: "development"}
	if devCfg.IsProduction() {
		t.Error("Expected IsProduction() to return false for development env")
	}
}

// TestConfig_IsDevelopment verifies development detection
func TestConfig_IsDevelopment(t *testing.T) {
	devCfg := &Config{Env: "development"}
	if !devCfg.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return true for development env")
	}

	prodCfg := &Config{Env: "production"}
	if prodCfg.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return false for production env")
	}
}

// TestConfig_Validate_MultipleErrors verifies multiple validation errors are reported
func TestConfig_Validate_MultipleErrors(t *testing.T) {
	cfg := &Config{
		Port:              "invalid",
		Env:               "invalid",
		SessionSecret:     "",
		RootAdminPassword: "",
		CookieHttpOnly:    true,
		CSRFEnabled:       true,
		CookieSameSite:    "invalid",
		SessionMaxAge:     -1,
		LogLevel:          "invalid",
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("Expected validation to fail")
	}

	errStr := err.Error()
	requiredErrors := []string{"SESSION_SECRET", "ROOT_ADMIN_PASSWORD", "PORT", "ENV", "SAMESITE", "MAX_AGE", "LOG_LEVEL"}

	for _, required := range requiredErrors {
		if !contains(errStr, required) {
			t.Errorf("Expected error to mention %s, but it wasn't found", required)
		}
	}
}

// TestEnsureDir verifies directory creation
func TestEnsureDir(t *testing.T) {
	testDir := "./test_uploads_" + strconv.Itoa(int(os.Getpid()))

	// Clean up after test
	defer os.RemoveAll(testDir)

	// Test creating directory
	err := ensureDir(testDir)
	if err != nil {
		t.Errorf("Expected ensureDir to succeed, got error: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Expected directory to be created")
	}

	// Test creating nested directory
	nestedDir := testDir + "/nested/path"
	err = ensureDir(nestedDir)
	if err != nil {
		t.Errorf("Expected ensureDir to create nested path, got error: %v", err)
	}
}

// TestLogLevelCaseInsensitive verifies log level is case insensitive
func TestLogLevelCaseInsensitive(t *testing.T) {
	levels := []string{"DEBUG", "INFO", "WARN", "ERROR", "Debug", "Info"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			clearEnvVars()
			os.Setenv("LOG_LEVEL", level)
			cfg := Load()

			// The loaded value should be lowercased - just verify it's a valid level
			validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
			if !validLevels[cfg.LogLevel] {
				t.Errorf("Expected log level to be valid, got: %s", cfg.LogLevel)
			}
		})
	}
}

// helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// clearEnvVars clears all configuration environment variables for clean testing
func clearEnvVars() {
	vars := []string{
		"PORT", "ENV", "DATABASE_URL", "SESSION_SECRET", "SESSION_MAX_AGE",
		"COOKIE_SECURE", "COOKIE_HTTPONLY", "COOKIE_SAMESITE", "CSRF_ENABLED",
		"TRUSTED_PROXIES", "ROOT_ADMIN_USERNAME", "ROOT_ADMIN_PASSWORD",
		"UPLOAD_PATH", "MAX_UPLOAD_SIZE", "LOG_LEVEL",
	}
	for _, v := range vars {
		os.Unsetenv(v)
	}
}
