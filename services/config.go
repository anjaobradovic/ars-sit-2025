package services

import (
	"errors"
	"log"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/repositories"
	"github.com/google/uuid"
)

type ConfigService struct {
	repo *repositories.ConfigRepository
}

func NewConfigService(repo *repositories.ConfigRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

func (s *ConfigService) Create(config *model.Config) error {
	log.Printf("Before UUID: %+v\n", config) // <--- dodaj ovo
	if config.ID == "" {
		config.ID = uuid.New().String()
	}
	log.Printf("After UUID: %+v\n", config) // <--- i ovo
	if config.Name == "" {
		return errors.New("name is required")
	}
	if config.Version == "" {
		return errors.New("version is required")
	}
	return s.repo.Save(*config)
}

func (s *ConfigService) Get(name, version string) (*model.Config, error) {
	if name == "" || version == "" {
		return nil, errors.New("name and version are required")
	}
	return s.repo.GetByNameAndVersion(name, version)
}

// Delete by ID + version
func (s *ConfigService) Delete(name, version string) error {
	if name == "" || version == "" {
		return errors.New("name and version are required")
	}
	return s.repo.DeleteByNameAndVersion(name, version)
}
