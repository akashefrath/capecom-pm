package v1

import (
	"capecom-pm/internal/cache"
	"capecom-pm/internal/container"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(v1 *gin.RouterGroup, c *container.Container) {
	h := c.Handler.AuthHandler
	auth := v1.Group("/auth")
	auth.Use(cache.SlidingWindowRateLimiter(c.RedisClient, 30, time.Minute))
	auth.POST("/login", h.Login)
	auth.POST("/refresh-token", h.Refresh)
	auth.GET("/me", c.Middleware.AuthMiddleware.VerifyToken(), h.Me)
	auth.POST("/logout", c.Middleware.AuthMiddleware.VerifyToken(), h.Logout)
}
