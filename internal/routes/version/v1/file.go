package v1

import (
	"capecom-pm/internal/container"

	"github.com/gin-gonic/gin"
)

func FileRoutes(v1 *gin.RouterGroup, c *container.Container) {
	h := c.Handler.FileHandler
	file := v1.Group("/files")
	file.Use(c.Middleware.AuthMiddleware.VerifyToken())
	file.POST("", h.CreateFile)
	file.POST("/confirm-upload", h.ConfirmUpload)
}
