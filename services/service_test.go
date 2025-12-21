package services

import (
	"context"
	"testing"

	"github.com/anjaobradovic/ars-sit-2025/model"
)

func TestCreateConfig_MissingName(t *testing.T) {
	service := NewConfigService(nil)

	cfg := &model.Config{
		Version: "1.0",
	}

	err := service.Create(context.Background(), cfg)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateConfig_MissingVersion(t *testing.T) {
	service := NewConfigService(nil)

	cfg := &model.Config{
		Name: "test",
	}

	err := service.Create(context.Background(), cfg)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
