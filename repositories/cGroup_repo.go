package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/hashicorp/consul/api"
)

type GroupRepository struct {
	kv *api.KV
}

func NewGroupRepository(consulAddr string) (*GroupRepository, error) {
	cfg := api.DefaultConfig()
	cfg.Address = consulAddr

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &GroupRepository{kv: client.KV()}, nil
}

func groupKey(name, version string) string {
	return fmt.Sprintf("groups/%s/%s", name, version)
}

func (r *GroupRepository) Save(group model.ConfigurationGroup) error {
	key := groupKey(group.Name, group.Version)

	existing, _, err := r.kv.Get(key, nil)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("group with this name and version already exists")
	}

	data, err := json.Marshal(group)
	if err != nil {
		return err
	}

	log.Printf("Repository: saving new group %s %s", group.Name, group.Version)
	_, err = r.kv.Put(&api.KVPair{
		Key:   key,
		Value: data,
	}, nil)

	return err
}

func (r *GroupRepository) GetByNameAndVersion(name, version string) (*model.ConfigurationGroup, error) {
	key := groupKey(name, version)

	pair, _, err := r.kv.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, errors.New("group not found")
	}

	var group model.ConfigurationGroup
	if err := json.Unmarshal(pair.Value, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *GroupRepository) DeleteByNameAndVersion(name, version string) error {
	key := groupKey(name, version)
	_, err := r.kv.Delete(key, nil)
	return err
}

func (r *GroupRepository) Update(group model.ConfigurationGroup) error {
	key := groupKey(group.Name, group.Version)

	data, err := json.Marshal(group)
	if err != nil {
		return err
	}

	log.Printf("Repository: updating group %s %s with configs: %+v", group.Name, group.Version, group.Configurations)

	_, err = r.kv.Put(&api.KVPair{
		Key:   key,
		Value: data,
	}, nil)

	return err
}
