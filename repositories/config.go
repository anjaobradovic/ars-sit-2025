package repositories

import (
	"errors"
	"sync"

	"github.com/anjaobradovic/ars-sit-2025/model"
)

var (
	configs = make(map[string]model.Config)
	mu      sync.Mutex
)

func Save(config model.Config) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := configs[config.ID]; exists {
		return errors.New("config already exists")
	}

	configs[config.ID] = config
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
