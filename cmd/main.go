package main

import (
	"fmt"
	"log"

	"github.com/akashefrath/capecom-pm/internal/cache"
	"github.com/akashefrath/capecom-pm/internal/config"
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/akashefrath/capecom-pm/internal/database/migration"
	"github.com/akashefrath/capecom-pm/internal/database/seeds"
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
	redis := cache.NewRedis(appConfig.RedisAddress)
	if redis == nil {
		panic("failed to connect to redis")

	}

	migration.Migrate(db)
	var err = seeds.SeedMasterData(db)
	if err != nil {
		log.Fatal("failed to seed master data")

	}
	return container.New(db, &appConfig, redis)

}
