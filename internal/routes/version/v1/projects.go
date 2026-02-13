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

	// tickets
	tk := c.Handler.TicketHandler
	project.POST("/:projectId/tickets", tk.CreateTicket)
	project.GET("/:projectId/tickets", tk.GetTickets)
	project.GET("/:projectId/tickets/:ticketId", tk.GetTicket)
	project.PUT("/:projectId/tickets/:ticketId", tk.UpdateTicket)
	project.PATCH("/:projectId/tickets/:ticketId/lifecycle", tk.UpdateLifecycleStatus)
	project.PATCH("/:projectId/tickets/:ticketId/assignee", tk.UpdateAssignee)
	project.DELETE("/:projectId/tickets/:ticketId", tk.DeleteTicket)

	// time entries
	te := c.Handler.TimeEntryHandler
	project.POST("/:projectId/tickets/:ticketId/time-entries", te.CreateTimeEntry)
	project.GET("/:projectId/tickets/:ticketId/time-entries", te.GetTimeEntries)
	project.GET("/:projectId/tickets/:ticketId/time-entries/:entryId", te.GetTimeEntry)
	project.PUT("/:projectId/tickets/:ticketId/time-entries/:entryId", te.UpdateTimeEntry)
	project.DELETE("/:projectId/tickets/:ticketId/time-entries/:entryId", te.DeleteTimeEntry)

	// ticket history
	hh := c.Handler.HistoryHandler
	project.GET("/:projectId/tickets/:ticketId/history", hh.GetTicketHistory)
}
