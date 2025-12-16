package services

import (
	"github.com/anjaobradovic/ars-sit-2025/model"
	"github.com/anjaobradovic/ars-sit-2025/repositories"
)

func CreateConfig(c model.Config) error {
	// koristi AddConfig iz repositories
	return repositories.Save(c)
}
