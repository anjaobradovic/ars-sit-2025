package services

import (
	"errors"
	"log"

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
	if group.Name == "" {
		return errors.New("name is required")
	}
	if group.Version == "" {
		return errors.New("version is required")
	}

	if group.Id == "" {
		group.Id = uuid.New().String()
	}

	if group.Configurations == nil {
		group.Configurations = []*model.LabeledConfiguration{}
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

	if cfg.Configuration == nil {
		return errors.New("configuration field is required")
	}

	// Generiši ID-jeve ako nedostaju
	if cfg.Id == "" {
		cfg.Id = uuid.New().String()
	}
	if cfg.Configuration.ID == "" {
		cfg.Configuration.ID = uuid.New().String()
	}

	// Provera duplikata po NAME + VERSION
	for _, c := range group.Configurations {
		if c.Configuration.Name == cfg.Configuration.Name &&
			c.Configuration.Version == cfg.Configuration.Version {
			return errors.New("configuration already exists in group")
		}
	}

	group.Configurations = append(group.Configurations, &cfg)

	log.Printf("Service: added config %+v to group %s %s", cfg, name, version)

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
