package adminV1Routes

import (
	"capecom-pm/internal/container"

	"github.com/gin-gonic/gin"
)

func ClientRoutes(v1 *gin.RouterGroup, c *container.Container) {
	h := c.Handler.ClientHandler
	client := v1.Group("/clients")
	client.POST("", h.CreateClient)
	client.GET("/:id", h.GetClientByID)
	client.GET("", h.GetClients)
	client.PUT("/:id", h.UpdateClient)
}
