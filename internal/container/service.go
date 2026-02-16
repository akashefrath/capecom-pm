package container

import "github.com/akashefrath/capecom-pm/internal/src/service"

type Service struct {
	Auth *service.Auth
}

func NewService(repo *Repositories) *Service {
	return &Service{
		Auth: service.NewAuth(repo.Auth, repo.User),
	}
}
