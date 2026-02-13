package container

import (
	"capecom-pm/internal/handlers"
	projecthandler "capecom-pm/internal/handlers/project"
)

type Handler struct {
	AuthHandler         *handlers.AuthHandler
	UserHandler         *handlers.UserHandler
	FileHandler         *handlers.FileHandler
	ClientHandler       *handlers.ClientHandler
	ProjectHandler      *projecthandler.ProjectHandler
	ProjectAssetHandler *projecthandler.AssetHandler
	ProjectTeamHandler  *projecthandler.TeamHandler
	UtilsHandler        *handlers.UtilsHandler
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
		ProjectHandler: projecthandler.NewProjectHandler(
			service.ProjectService),
		ProjectAssetHandler: projecthandler.NewAssetHandler(
			service.ProjectAssetService),
		ProjectTeamHandler: projecthandler.NewTeamHandler(
			service.ProjectTeamService),
		UtilsHandler: handlers.NewUtilsHandler(
			service.UtilsService),
	}
}
