package container

import (
	"capecom-pm/internal/repositories"

	"gorm.io/gorm"
)

type Repository struct {
	AuthRepo *repositories.AuthRepo
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AuthRepo: repositories.NewAuthRepo(db),
	}
}
