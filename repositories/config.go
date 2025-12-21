package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/hashicorp/consul/api"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("repositories/config")

type ConfigRepository struct {
	kv *api.KV
}

func NewConfigRepository(addr string) (*ConfigRepository, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ConfigRepository{kv: client.KV()}, nil
}

func (r *ConfigRepository) Save(ctx context.Context, config model.Config) error {
	ctx, span := tracer.Start(ctx, "ConfigRepository.Save")
	defer span.End()

	key := fmt.Sprintf("configs/%s/%s", config.Name, config.Version)
	span.SetAttributes(
		attribute.String("consul.key", key),
		attribute.String("config.name", config.Name),
		attribute.String("config.version", config.Version),
	)

	// Optional: check existing
	{
		_, s := tracer.Start(ctx, "consul.kv.get")
		s.SetAttributes(attribute.String("consul.key", key))
		existing, _, err := r.kv.Get(key, nil)
		if err != nil {
			s.RecordError(err)
			s.SetStatus(codes.Error, "consul get failed")
			s.End()
			return err
		}
		s.End()

		if existing != nil {
			err := errors.New("configuration already exists")
			span.RecordError(err)
			span.SetStatus(codes.Error, "conflict")
			return err
		}
	}

	b, err := json.Marshal(config)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "marshal failed")
		return err
	}

	{
		_, s := tracer.Start(ctx, "consul.kv.put")
		s.SetAttributes(attribute.String("consul.key", key))
		_, err = r.kv.Put(&api.KVPair{Key: key, Value: b}, nil)
		if err != nil {
			s.RecordError(err)
			s.SetStatus(codes.Error, "consul put failed")
			s.End()
			return err
		}
		s.End()
	}

	return nil
}

func (r *ConfigRepository) GetByNameAndVersion(ctx context.Context, name, version string) (*model.Config, error) {
	ctx, span := tracer.Start(ctx, "ConfigRepository.GetByNameAndVersion")
	defer span.End()

	key := fmt.Sprintf("configs/%s/%s", name, version)
	span.SetAttributes(
		attribute.String("consul.key", key),
		attribute.String("config.name", name),
		attribute.String("config.version", version),
	)

	var pair *api.KVPair
	{
		_, s := tracer.Start(ctx, "consul.kv.get")
		s.SetAttributes(attribute.String("consul.key", key))
		var err error
		pair, _, err = r.kv.Get(key, nil)
		if err != nil {
			s.RecordError(err)
			s.SetStatus(codes.Error, "consul get failed")
			s.End()
			return nil, err
		}
		s.End()
	}

	if pair == nil {
		err := errors.New("configuration not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, "not found")
		return nil, err
	}

	var cfg model.Config
	if err := json.Unmarshal(pair.Value, &cfg); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unmarshal failed")
		return nil, err
	}

	return &cfg, nil
}

func (r *ConfigRepository) DeleteByNameAndVersion(ctx context.Context, name, version string) error {
	ctx, span := tracer.Start(ctx, "ConfigRepository.DeleteByNameAndVersion")
	defer span.End()

	key := fmt.Sprintf("configs/%s/%s", name, version)
	span.SetAttributes(
		attribute.String("consul.key", key),
		attribute.String("config.name", name),
		attribute.String("config.version", version),
	)

	{
		_, s := tracer.Start(ctx, "consul.kv.delete")
		s.SetAttributes(attribute.String("consul.key", key))
		_, err := r.kv.Delete(key, nil)
		if err != nil {
			s.RecordError(err)
			s.SetStatus(codes.Error, "consul delete failed")
			s.End()
			return err
		}
		s.End()
	}

	return nil
}
