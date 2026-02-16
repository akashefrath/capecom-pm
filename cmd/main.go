package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/akashefrath/capecom-pm/internal/config"
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/akashefrath/capecom-pm/internal/routes"
	"github.com/gofiber/fiber/v3"
)

func main() {
	appContainer := initApp()
	app := fiber.New(fiber.Config{

		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})
	routes.Setup(app, appContainer)
	err := app.Listen(fmt.Sprintf(":%d", appContainer.Config.Port))
	if err != nil {
		log.Fatal(err)
	}
}

func initApp() container.Container {
	appConfig := config.LoadEnv()
	db := config.InitDB(appConfig)
	return container.New(db, &appConfig)

}
