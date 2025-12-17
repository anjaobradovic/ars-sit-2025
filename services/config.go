package services

import (
	"errors"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/repositories"
)

func CreateConfig(c model.Config) error {
	// koristi AddConfig iz repositories
	return repositories.Save(c)
}

var configRepo = repositories.NewConfigRepository()

func AddConfig(cfg model.Config) error {

	// minimalna validacija
	if cfg.Name == "" || cfg.Version == "" {
		return errors.New("name and version are required")
	}

	// prosledjujemo repository-ju
	return configRepo.Add(cfg)
}

func GetConfig(name, version string) (model.Config, error) {
	if name == "" || version == "" {
		return model.Config{}, errors.New("name and version are required")
	}
	return configRepo.Get(name, version)
}
