package services

import (
	"errors"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/repositories"
	"github.com/google/uuid"
)

type GroupService struct {
	repo *repositories.GroupRepository
}

func NewGroupService(repo *repositories.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) Create(group *model.ConfigurationGroup) error {
	if group.Id == "" {
		group.Id = uuid.New().String()
	}
	if group.Name == "" {
		return errors.New("name is required")
	}
	if group.Version == "" {
		return errors.New("version is required")
	}

	// Ovo briše konfiguracije, što je problem:
	// group.Configurations = []*model.LabeledConfiguration{}

	// Umesto toga, možeš proći kroz svaku konfiguraciju i dodati UUID ako nedostaje:
	for _, cfg := range group.Configurations {
		if cfg.Id == "" {
			cfg.Id = uuid.New().String()
		}
		if cfg.Configuration != nil && cfg.Configuration.ID == "" {
			cfg.Configuration.ID = uuid.New().String()
		}
	}

	return s.repo.Save(*group)
}

func (s *GroupService) Get(name, version string) (*model.ConfigurationGroup, error) {
	if name == "" || version == "" {
		return nil, errors.New("name and version are required")
	}
	return s.repo.GetByNameAndVersion(name, version)
}

func (s *GroupService) Delete(name, version string) error {
	if name == "" || version == "" {
		return errors.New("name and version are required")
	}
	return s.repo.DeleteByNameAndVersion(name, version)
}

func (s *GroupService) AddConfig(name, version string, cfg model.LabeledConfiguration) error {
	group, err := s.repo.GetByNameAndVersion(name, version)
	if err != nil {
		return err
	}

	for _, c := range group.Configurations {
		if c.Id == cfg.Id {
			return errors.New("config already exists in group")
		}
	}

	group.Configurations = append(group.Configurations, &cfg)
	return s.repo.Update(*group)
}

func (s *GroupService) RemoveConfig(name, version, configID string) error {
	group, err := s.repo.GetByNameAndVersion(name, version)
	if err != nil {
		return err
	}

	filtered := []*model.LabeledConfiguration{}
	for _, c := range group.Configurations {
		if c.Id != configID {
			filtered = append(filtered, c)
		}
	}

	group.Configurations = filtered
	return s.repo.Update(*group)
}
