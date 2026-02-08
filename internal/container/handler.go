package container

import "capecom-pm/internal/handlers"

type Handler struct {
	AuthHandler *handlers.AuthHandler
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		AuthHandler: handlers.NewAuthHandler(
			service.AuthService),
	}
}
