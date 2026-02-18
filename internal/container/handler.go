package container

import (
	"github.com/akashefrath/capecom-pm/internal/src/handler"
	utilshandler "github.com/akashefrath/capecom-pm/internal/src/handler/utils"
)

type Handler struct {
	AuthHandler  handler.AuthHandler
	UserHandler  handler.UserHandler
	RoleHandler  handler.RoleHandler
	UtilsHandler utilshandler.UtilsHandler
}

func SetupHandler(service *Service) *Handler {
	return &Handler{
		AuthHandler:  handler.NewAuth(service.Auth),
		UserHandler:  handler.NewUser(service.User),
		RoleHandler:  handler.NewRole(service.Role),
		UtilsHandler: utilshandler.NewUtilsHandler(service.Utils),
	}
}
