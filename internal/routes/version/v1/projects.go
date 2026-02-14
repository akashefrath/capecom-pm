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
	projectNs := v1.Group("/project")
	projectNs.Use(c.Middleware.AuthMiddleware.VerifyToken())
	project.Use(c.Middleware.AuthMiddleware.VerifyToken())
	project.Use(c.Middleware.RABCMiddleware.IsManagerOrAdmin())
	project.POST("", h.CreateProject)
	projectNs.GET("", h.GetProjects)           // for all with conditions
	projectNs.GET("/:projectId", h.GetProject) // for all with conditions
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
	projectNs.GET("/:projectId/tickets", tk.GetTickets)          // for all with conditions
	projectNs.GET("/:projectId/tickets/:ticketId", tk.GetTicket) // for all with conditions
	project.PUT("/:projectId/tickets/:ticketId", tk.UpdateTicket)
	project.PATCH("/:projectId/tickets/:ticketId/lifecycle", tk.UpdateLifecycleStatus)
	project.PATCH("/:projectId/tickets/:ticketId/assignee", tk.UpdateAssignee)
	project.DELETE("/:projectId/tickets/:ticketId", tk.DeleteTicket)

	// time entries
	te := c.Handler.TimeEntryHandler
	projectNs.POST("/:projectId/tickets/:ticketId/time-entries", te.CreateTimeEntry)            // for all with conditions
	projectNs.GET("/:projectId/tickets/:ticketId/time-entries", te.GetTimeEntries)              // for all with conditions
	projectNs.GET("/:projectId/tickets/:ticketId/time-entries/:entryId", te.GetTimeEntry)       // for all with conditions
	projectNs.PUT("/:projectId/tickets/:ticketId/time-entries/:entryId", te.UpdateTimeEntry)    // for all with conditions
	projectNs.DELETE("/:projectId/tickets/:ticketId/time-entries/:entryId", te.DeleteTimeEntry) // for all with conditions

	// ticket history
	hh := c.Handler.HistoryHandler
	projectNs.GET("/:projectId/tickets/:ticketId/history", hh.GetTicketHistory)
}
