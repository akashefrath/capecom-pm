package container

import (
	"capecom-pm/internal/services"
	jwtutil "capecom-pm/internal/utils/jwt"
)

type Service struct {
	AuthService *services.AuthService
	UserService *services.UserService
}

func NewService(repository *Repository, jwt *jwtutil.Manager) *Service {
	return &Service{
		AuthService: services.NewAuthService(repository.AuthRepo, jwt, repository.UserRepo, repository.SessionRepo, repository.CacheRepo),
		UserService: services.NewUserService(
			repository.UserRepo,
			repository.MasterDataRepo,
		),
	}
}
