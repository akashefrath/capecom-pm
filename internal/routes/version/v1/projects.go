package v1

import (
	"capecom-pm/internal/container"

	"github.com/gin-gonic/gin"
)

func Projects(v1 *gin.RouterGroup, c *container.Container) {
	h := c.Handler.ProjectHandler
	ah := c.Handler.ProjectAssetHandler
	th := c.Handler.ProjectTeamHandler

	project := v1.Group("/project")
	project.Use(c.Middleware.AuthMiddleware.VerifyToken())
	project.Use(c.Middleware.RABCMiddleware.IsManagerOrAdmin())
	project.POST("", h.CreateProject)
	project.GET("", h.GetProjects)
	project.PUT("/update/:id", h.UpdateProject)
	project.PATCH("/lifecycle/:id", h.UpdateLifecycleStatus)
	project.DELETE("/delete/:id", h.DeleteProject)

	// project assets
	project.POST("/:projectId/assets", ah.CreateAsset)
	project.GET("/:projectId/assets", ah.GetAssets)
	project.GET("/:projectId/assets/:assetId", ah.GetAsset)
	project.PUT("/:projectId/assets/:assetId", ah.UpdateAsset)
	project.DELETE("/:projectId/assets/:assetId", ah.DeleteAsset)

	// project managers
	project.POST("/:projectId/managers", th.AssignManagers)
	project.DELETE("/:projectId/managers", th.RemoveManagers)
	project.GET("/:projectId/managers", th.GetManagers)

	// project members
	project.POST("/:projectId/members", th.AssignMembers)
	project.DELETE("/:projectId/members", th.RemoveMembers)
	project.GET("/:projectId/members", th.GetMembers)
}
