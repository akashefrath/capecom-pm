package container

import "github.com/akashefrath/capecom-pm/internal/src/handler"

type Handler struct {
	AuthHandler handler.AuthHandler
}

func SetupHandler(service *Service) *Handler {
	return &Handler{
		AuthHandler: handler.NewAuth(service.Auth),
	}
}
