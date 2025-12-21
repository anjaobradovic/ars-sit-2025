package model

import "testing"

func TestConfigCreation(t *testing.T) {
	cfg := Config{
		Name:    "test-config",
		Version: "1.0",
		Parameters: map[string]string{
			"db.host": "localhost",
		},
	}

	if cfg.Name != "test-config" {
		t.Errorf("expected name 'test-config', got %s", cfg.Name)
	}

	if cfg.Version != "1.0" {
		t.Errorf("expected version '1.0', got %s", cfg.Version)
	}

	if cfg.Parameters["db.host"] != "localhost" {
		t.Errorf("expected db.host to be localhost")
	}
}
