package services

import (
	"context"
	"errors"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/repositories"
	"github.com/google/uuid"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("services/config")

type ConfigService struct {
	repo *repositories.ConfigRepository
}

func NewConfigService(repo *repositories.ConfigRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

func (s *ConfigService) Create(ctx context.Context, config *model.Config) error {
	ctx, span := tracer.Start(ctx, "ConfigService.Create")
	defer span.End()

	span.SetAttributes(
		attribute.String("config.name", config.Name),
		attribute.String("config.version", config.Version),
	)

	if config.Name == "" || config.Version == "" {
		err := errors.New("name and version are required")
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return err
	}

	config.ID = uuid.NewString()

	if err := s.repo.Save(ctx, *config); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "repo save failed")
		return err
	}

	return nil
}

func (s *ConfigService) Get(ctx context.Context, name, version string) (*model.Config, error) {
	ctx, span := tracer.Start(ctx, "ConfigService.Get")
	defer span.End()

	span.SetAttributes(
		attribute.String("config.name", name),
		attribute.String("config.version", version),
	)

	cfg, err := s.repo.GetByNameAndVersion(ctx, name, version)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "repo get failed")
		return nil, err
	}

	return cfg, nil
}

func (s *ConfigService) Delete(ctx context.Context, name, version string) error {
	ctx, span := tracer.Start(ctx, "ConfigService.Delete")
	defer span.End()

	span.SetAttributes(
		attribute.String("config.name", name),
		attribute.String("config.version", version),
	)

	if err := s.repo.DeleteByNameAndVersion(ctx, name, version); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "repo delete failed")
		return err
	}

	return nil
}
