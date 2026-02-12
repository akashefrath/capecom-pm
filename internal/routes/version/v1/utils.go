package v1

import (
	"capecom-pm/internal/container"

	"github.com/gin-gonic/gin"
)

func UtilsRoutes(v1 *gin.RouterGroup, c *container.Container) {
	h := c.Handler.UtilsHandler
	utils := v1.Group("/utils")

	utils.GET("", h.GetAll)
	utils.GET("/roles", h.GetRoles)
	utils.GET("/user-groups", h.GetUserGroups)
	utils.GET("/designations", h.GetDesignations)
	utils.GET("/departments", h.GetDepartments)
	utils.GET("/clients", h.GetClients)
	utils.GET("/ticket-types", h.GetTicketTypes)
	utils.GET("/users", h.GetUsers)
}
