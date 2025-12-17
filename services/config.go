package services

import (
	"errors"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/repositories"
)

type ConfigService struct {
	repo *repositories.ConfigRepository
}

func NewConfigService(repo *repositories.ConfigRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

func (s *ConfigService) Create(config model.Config) error {
	if config.ID == "" {
		return errors.New("id is required")
	}
	if config.Name == "" {
		return errors.New("name is required")
	}
	if config.Version == "" {
		return errors.New("version is required")
	}

	return s.repo.Save(config)
}
