package adminV1Routes

import (
	"capecom-pm/internal/container"

	"github.com/gin-gonic/gin"
)

func UserRoutes(v1 *gin.RouterGroup, c *container.Container) {
	h := c.Handler.UserHandler
	user := v1.Group("/users")
	user.POST("", h.CreateUser)

}
