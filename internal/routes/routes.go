package routes

import (
	"capecom-pm/internal/cache"
	"capecom-pm/internal/container"
	"capecom-pm/internal/routes/version/v1"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, c *container.Container) {
	r.Any("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})

	apiV1 := r.Group("/api/v1")
	apiV1.Use(cache.SlidingWindowRateLimiter(c.RedisClient, 100, time.Minute))
	v1.Routes(apiV1, c)

}
