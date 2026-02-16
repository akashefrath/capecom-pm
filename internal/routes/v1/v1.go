package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App, container container.Container) {
	v1API := app.Group("api/v1")
	Auth(v1API, container)

}
