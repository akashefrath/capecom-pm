package v1

import (
	"capecom-pm/internal/container"

	"github.com/gin-gonic/gin"
)

func Projects(v1 *gin.RouterGroup, c *container.Container) {
	h := c.Handler.ProjectHandler
	project := v1.Group("/project")
	project.Use(c.Middleware.AuthMiddleware.VerifyToken())
	project.Use(c.Middleware.RABCMiddleware.IsManagerOrAdmin())
	project.POST("", h.CreateProject)
}
