-- Lab Settings table for storing lab configuration as key-value pairs
-- Used for lab name, description, and other configurable settings
-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- Lab settings table: stores configuration as key-value pairs
-- All settings are unique by key
CREATE TABLE lab_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    setting_key TEXT NOT NULL,
    setting_value TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Unique constraint on setting_key to prevent duplicates
CREATE UNIQUE INDEX idx_lab_settings_key ON lab_settings(setting_key);

-- Insert default lab settings
INSERT INTO lab_settings (setting_key, setting_value) VALUES
    ('lab_name', 'Research Lab'),
    ('lab_description', 'A research laboratory');
