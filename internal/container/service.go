package container

import (
	"capecom-pm/internal/services"

	"gorm.io/gorm"
)

type Service struct {
	AuthService *services.AuthService
}

func NewService(db *gorm.DB, repository *Repository) *Service {
	return &Service{
		AuthService: services.NewAuthService(repository.AuthRepo),
	}
}
