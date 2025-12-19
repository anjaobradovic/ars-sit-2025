package repositories

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/hashicorp/consul/api"
)

type ConfigRepository struct {
	kv *api.KV
}

func NewConfigRepository(consulAddr string) (*ConfigRepository, error) {
	cfg := api.DefaultConfig()
	cfg.Address = consulAddr

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ConfigRepository{kv: client.KV()}, nil
}

// Save a new config (immutable, idempotent)
func (r *ConfigRepository) Save(config model.Config) error {
	key := fmt.Sprintf("configs/%s/%s", config.ID, config.Version)

	existing, _, err := r.kv.Get(key, nil)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("config already exists")
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = r.kv.Put(&api.KVPair{
		Key:   key,
		Value: data,
	}, nil)

	return err
}

// Get config by ID + version
func (r *ConfigRepository) GetByIDAndVersion(id, version string) (*model.Config, error) {
	key := fmt.Sprintf("configs/%s/%s", id, version)
	pair, _, err := r.kv.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, errors.New("config not found")
	}

	var config model.Config
	if err := json.Unmarshal(pair.Value, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Delete config by ID + version
func (r *ConfigRepository) DeleteByIDAndVersion(id, version string) error {
	key := fmt.Sprintf("configs/%s/%s", id, version)
	_, err := r.kv.Delete(key, nil)
	return err
}
