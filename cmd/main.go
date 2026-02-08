package main

import (
	"capecom-pm/internal/bootstrap"
	"capecom-pm/internal/cache"
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
	var err = bootstrap.SeedMasterData(db)
	if err != nil {
		log.Fatal("failed to seed master data")
		return
	}
	redis := cache.NewRedis(appConfig.RedisAddress)
	if redis == nil {
		log.Fatal("failed to connect to redis")
		return
	}
	c := container.NewContainer(db, appConfig, redis)
	r := gin.Default()
	routes.Setup(r, c)
	fmt.Println("http://localhost:" + appConfig.Port)
	err = r.Run(":" + appConfig.Port)
	if err != nil {
		return
	}
}
