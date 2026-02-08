package v1

import (
	"capecom-pm/internal/container"

	"github.com/gin-gonic/gin"
)

func Routes(v1 *gin.RouterGroup, c *container.Container) {
	AuthRoutes(v1, c)

}
