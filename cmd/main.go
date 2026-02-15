package main

import (
	"capecom-pm/internal/bootstrap"
	"capecom-pm/internal/cache"
	"capecom-pm/internal/config"
	"capecom-pm/internal/container"
	"capecom-pm/internal/routes"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
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
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    cors.DefaultConfig().AllowMethods,
		AllowHeaders:    cors.DefaultConfig().AllowHeaders,
		ExposeHeaders:   []string{"Content-Length"},
		 
		MaxAge: 12 * time.Hour,
	}))

	routes.Setup(r, c)
	fmt.Println("http://localhost:" + appConfig.Port)
	err = r.Run(":" + appConfig.Port)
	if err != nil {
		return
	}
}
