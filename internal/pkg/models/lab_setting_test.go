package models

import (
	"testing"
	"time"
)

func TestLabSetting_New(t *testing.T) {
	setting := LabSetting{
		ID:           1,
		SettingKey:   LabSettingName,
		SettingValue: "Test Lab Name",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if setting.ID != 1 {
		t.Errorf("expected ID to be 1, got %d", setting.ID)
	}

	if setting.SettingKey != LabSettingName {
		t.Errorf("expected SettingKey to be %s, got %s", LabSettingName, setting.SettingKey)
	}

	if setting.SettingValue != "Test Lab Name" {
		t.Errorf("expected SettingValue to be 'Test Lab Name', got %s", setting.SettingValue)
	}
}

func TestLabSetting_Constants(t *testing.T) {
	if LabSettingName != "lab_name" {
		t.Errorf("expected LabSettingName to be 'lab_name', got %s", LabSettingName)
	}

	if LabSettingDescription != "lab_description" {
		t.Errorf("expected LabSettingDescription to be 'lab_description', got %s", LabSettingDescription)
	}
}

func TestLabSetting_FieldTags(t *testing.T) {
	setting := LabSetting{
		ID:           1,
		SettingKey:   "test_key",
		SettingValue: "test_value",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Test that all fields are settable
	if setting.ID == 0 {
		t.Error("ID should be settable")
	}

	if setting.SettingKey == "" {
		t.Error("SettingKey should be settable")
	}

	if setting.SettingValue == "" {
		t.Error("SettingValue should be settable")
	}

	if setting.CreatedAt.IsZero() {
		t.Error("CreatedAt should be settable")
	}

	if setting.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be settable")
	}
}
