package v1

import (
	"capecom-pm/internal/container"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(v1 *gin.RouterGroup, c *container.Container) {
	h := c.Handler.AuthHandler
	auth := v1.Group("/auth")

	auth.POST("/login", h.Login)
	auth.POST("/refresh-token", h.Refresh)
	auth.GET("/me", c.Middleware.UserMiddleware.VerifyUserToken(), h.Me)
}
