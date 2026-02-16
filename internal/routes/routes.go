package routes

import (
	"net/http"

	"github.com/akashefrath/capecom-pm/internal/container"
	v1 "github.com/akashefrath/capecom-pm/internal/routes/v1"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App, container container.Container) {
	//	app.Use(r.New())
	app.All("/ping", func(c fiber.Ctx) error {
		return response.JSON(c, http.StatusOK, response.APIResponse{
			Success: true,
			Message: "Pong",
		})
	})

	v1.Setup(app, container)
}
