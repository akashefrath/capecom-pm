package adminV1Routes

import (
	"capecom-pm/internal/container"

	"github.com/gin-gonic/gin"
)

func SetAdminRoute(v1 *gin.RouterGroup, c *container.Container) {

	admin := v1.Group("/admin")
	UserRoutes(admin, c)
}
