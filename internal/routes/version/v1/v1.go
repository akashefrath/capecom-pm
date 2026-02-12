package v1

import (
	"capecom-pm/internal/container"
	adminV1Routes "capecom-pm/internal/routes/version/v1/admin"

	"github.com/gin-gonic/gin"
)

func Routes(v1 *gin.RouterGroup, c *container.Container) {
	adminV1Routes.SetAdminRoute(v1, c)
	AuthRoutes(v1, c)
	Projects(v1, c)

}
