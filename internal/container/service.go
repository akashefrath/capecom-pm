package container

import (
	"github.com/akashefrath/capecom-pm/internal/src/service"
	jwtutil "github.com/akashefrath/capecom-pm/internal/utils/jwt"
)

type Service struct {
	Auth *service.Auth
}

func NewService(repo *Repositories, jwtManager *jwtutil.Manager) *Service {
	return &Service{
		Auth: service.NewAuth(repo.Session, repo.User, jwtManager),
	}
}
