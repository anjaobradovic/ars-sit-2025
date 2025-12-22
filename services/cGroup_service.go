package services

import (
	"errors"
	"log"
	"strings"

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

// parseLabels parses "k1:v1;k2:v2" into a map.
// Returns error if format is invalid.
func parseLabels(raw string) (map[string]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("labels query param is required")
	}

	out := map[string]string{}
	pairs := strings.Split(raw, ";")
	for _, p := range pairs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, ":", 2)
		if len(kv) != 2 {
			return nil, errors.New("invalid labels format, expected key:value;key2:value2")
		}

		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if k == "" || v == "" {
			return nil, errors.New("invalid labels format, empty key or value")
		}
		out[k] = v
	}

	if len(out) == 0 {
		return nil, errors.New("labels query param is required")
	}
	return out, nil
}

// matchesAllLabels checks AND matching: all query labels must exist and equal in cfg.Labels.
func matchesAllLabels(cfg *model.LabeledConfiguration, query map[string]string) bool {
	for k, v := range query {
		if cfg.Labels == nil || cfg.Labels[k] != v {
			return false
		}
	}
	return true
}

// DeleteConfigsByLabels removes all labeled configurations from a group that match ALL labels.
func (s *GroupService) DeleteConfigsByLabels(name, version, rawLabels string) (int, error) {
	if name == "" || version == "" {
		return 0, errors.New("name and version are required")
	}

	queryLabels, err := parseLabels(rawLabels)
	if err != nil {
		return 0, err
	}

	group, err := s.repo.GetByNameAndVersion(name, version)
	if err != nil {
		return 0, err
	}

	kept := make([]*model.LabeledConfiguration, 0, len(group.Configurations))
	deleted := 0

	for _, cfg := range group.Configurations {
		if matchesAllLabels(cfg, queryLabels) {
			deleted++
			continue
		}
		kept = append(kept, cfg)
	}

	group.Configurations = kept
	if err := s.repo.Update(*group); err != nil {
		return 0, err
	}

	log.Printf("Service: deleted %d configs by labels from group %s %s", deleted, name, version)
	return deleted, nil
}
