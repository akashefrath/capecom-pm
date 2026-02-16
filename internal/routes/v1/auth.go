package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/gofiber/fiber/v3"
)

func Auth(route fiber.Router, container container.Container) {
	auth := route.Group("auth")
	authHandler := container.Handler.AuthHandler
	auth.Post("login", authHandler.Login)

}
