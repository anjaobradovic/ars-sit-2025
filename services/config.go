package services

import (
	"errors"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/repositories"
)

type ConfigService struct {
	repo *repositories.ConfigRepository
}

// add
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

// found
func (s *ConfigService) Get(id string) (*model.Config, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(id)
}

// delete
func (s *ConfigService) Delete(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.DeleteByID(id)
}
