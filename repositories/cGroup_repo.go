package repositories

import (
	"encoding/json"
	"errors"
	"fmt"

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

func (r *GroupRepository) Save(group model.ConfigurationGroup) error {
	key := fmt.Sprintf("groups/%s", group.Id)

	existing, _, err := r.kv.Get(key, nil)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("group already exists")
	}

	data, err := json.Marshal(group)
	if err != nil {
		return err
	}

	_, err = r.kv.Put(&api.KVPair{
		Key:   key,
		Value: data,
	}, nil)
	return err
}

func (r *GroupRepository) GetByID(id string) (*model.ConfigurationGroup, error) {
	key := fmt.Sprintf("groups/%s", id)
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

func (r *GroupRepository) DeleteByID(id string) error {
	key := fmt.Sprintf("groups/%s", id)
	_, err := r.kv.Delete(key, nil)
	return err
}

func (r *GroupRepository) Update(group model.ConfigurationGroup) error {
	key := fmt.Sprintf("groups/%s", group.Id)

	existing, _, err := r.kv.Get(key, nil)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("group does not exist")
	}

	data, err := json.Marshal(group)
	if err != nil {
		return err
	}

	_, err = r.kv.Put(&api.KVPair{
		Key:   key,
		Value: data,
	}, nil)
	return err
}
