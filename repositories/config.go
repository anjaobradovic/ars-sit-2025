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
