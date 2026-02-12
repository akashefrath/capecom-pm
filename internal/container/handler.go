package container

import "capecom-pm/internal/handlers"

type Handler struct {
	AuthHandler    *handlers.AuthHandler
	UserHandler    *handlers.UserHandler
	FileHandler    *handlers.FileHandler
	ClientHandler  *handlers.ClientHandler
	ProjectHandler *handlers.ProjectHandler
	UtilsHandler   *handlers.UtilsHandler
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		AuthHandler: handlers.NewAuthHandler(
			service.AuthService),
		UserHandler: handlers.NewUserHandler(
			service.UserService),
		FileHandler: handlers.NewFileHandler(
			service.FileService),
		ClientHandler: handlers.NewClientHandler(
			service.ClientService),
		ProjectHandler: handlers.NewProjectHandler(
			service.ProjectService),
		UtilsHandler: handlers.NewUtilsHandler(
			service.UtilsService),
	}
}
