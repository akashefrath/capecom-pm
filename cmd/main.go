package main

import (
	"fmt"

	"github.com/akashefrath/capecom-pm/internal/config"
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/akashefrath/capecom-pm/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	appContainer := initApp()
	r := gin.Default()

	routes.Setup(r, appContainer)
	err := r.Run(fmt.Sprintf(":%d", appContainer.Config.Port))
	if err != nil {
		return
	}
}

func initApp() container.Container {
	appConfig := config.LoadEnv()
	db := config.InitDB(appConfig)
	return container.New(db, &appConfig)

}
