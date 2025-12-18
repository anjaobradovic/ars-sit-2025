package services

import (
	"errors"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/repositories"
)

type GroupService struct {
	repo *repositories.GroupRepository
}

func NewGroupService(repo *repositories.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) Create(group model.ConfigurationGroup) error {
	if group.Id == "" {
		return errors.New("id is required")
	}
	if group.Name == "" {
		return errors.New("name is required")
	}
	if group.Version == "" {
		return errors.New("version is required")
	}
	return s.repo.Save(group)
}

func (s *GroupService) Get(id string) (*model.ConfigurationGroup, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(id)
}

func (s *GroupService) Delete(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.DeleteByID(id)
}

func (s *GroupService) AddConfigToGroup(groupID string, labeledConfig model.LabeledConfiguration) error {
	group, err := s.repo.GetByID(groupID)
	if err != nil {
		return err
	}

	for _, c := range group.Configurations {
		if c.Id == labeledConfig.Id {
			return errors.New("config already in group")
		}
	}

	group.Configurations = append(group.Configurations, &labeledConfig)
	return s.repo.Update(*group)
}

func (s *GroupService) RemoveConfigFromGroup(groupID string, labeledConfigID string) error {
	group, err := s.repo.GetByID(groupID)
	if err != nil {
		return err
	}

	newConfigs := []*model.LabeledConfiguration{}
	for _, c := range group.Configurations {
		if c.Id != labeledConfigID {
			newConfigs = append(newConfigs, c)
		}
	}

	group.Configurations = newConfigs
	return s.repo.Update(*group)
}
