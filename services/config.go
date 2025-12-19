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

// Get by ID + version
func (s *ConfigService) Get(id, version string) (*model.Config, error) {
	if id == "" || version == "" {
		return nil, errors.New("id and version are required")
	}
	return s.repo.GetByIDAndVersion(id, version)
}

// Delete by ID + version
func (s *ConfigService) Delete(id, version string) error {
	if id == "" || version == "" {
		return errors.New("id and version are required")
	}
	return s.repo.DeleteByIDAndVersion(id, version)
}
