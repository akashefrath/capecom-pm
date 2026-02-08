package container

import (
	"capecom-pm/internal/services"
	jwtutil "capecom-pm/internal/utils/jwt"

	"gorm.io/gorm"
)

type Service struct {
	AuthService *services.AuthService
	UserService *services.UserService
}

func NewService(db *gorm.DB, repository *Repository, jwt *jwtutil.JWTManager) *Service {
	return &Service{
		AuthService: services.NewAuthService(repository.AuthRepo, jwt),
		UserService: services.NewUserService(
			repository.UserRepo,
			repository.MasterDataRepo,
		),
	}
}
