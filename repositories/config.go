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

// add new config
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

// found config by id
func (r *ConfigRepository) GetByID(id string) (*model.Config, error) {
	key := fmt.Sprintf("configs/%s", id)
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

// delete config by id
func (r *ConfigRepository) DeleteByID(id string) error {
	key := fmt.Sprintf("configs/%s", id)
	_, err := r.kv.Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}

type ConfigRepository struct {
	data map[string]model.Config
}

func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{
		data: make(map[string]model.Config),
	}
}

func (r *ConfigRepository) Add(cfg model.Config) error {

	// kljuc = name:version
	key := cfg.Name + ":" + cfg.Version

	// ako vec postoji -> greska
	if _, exists := r.data[key]; exists {
		return errors.New("config already exists")
	}

	// cuvamo konfiguraciju
	r.data[key] = cfg
	return nil
}

func (r *ConfigRepository) Get(name, version string) (model.Config, error) {
	key := name + ":" + version

	cfg, exists := r.data[key]
	if !exists {
		return model.Config{}, errors.New("config not found")
	}
	return cfg, nil
}
