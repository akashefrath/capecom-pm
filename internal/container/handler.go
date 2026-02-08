package container

import "capecom-pm/internal/handlers"

type Handler struct {
	AuthHandler *handlers.AuthHandler
	UserHandler *handlers.UserHandler
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		AuthHandler: handlers.NewAuthHandler(
			service.AuthService),
		UserHandler: handlers.NewUserHandler(
			service.UserService),
	}
}
