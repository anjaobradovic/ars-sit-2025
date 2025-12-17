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

//kv -  key value

func NewConfigRepository(consulAddr string) (*ConfigRepository, error) {
	cfg := api.DefaultConfig()
	cfg.Address = consulAddr

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ConfigRepository{
		kv: client.KV(),
	}, nil
}

func (r *ConfigRepository) Save(config model.Config) error {
	key := fmt.Sprintf("configs/%s", config.ID)

	// check does it already exist
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
