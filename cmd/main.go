package main

import (
	"capecom-pm/internal/bootstrap"
	"capecom-pm/internal/config"
	"capecom-pm/internal/container"
	"capecom-pm/internal/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	appConfig := config.LoadEnv()
	db := config.ConnectDB(appConfig.DB)
	if db == nil {
		log.Fatal("failed to connect to database")
	}
	bootstrap.SeedMasterData(db)
	c := container.NewContainer(db, appConfig)
	r := gin.Default()
	routes.Setup(r, c)
	fmt.Println("http://localhost:" + appConfig.Port)
	err := r.Run(":" + appConfig.Port)
	if err != nil {
		return
	}
}
