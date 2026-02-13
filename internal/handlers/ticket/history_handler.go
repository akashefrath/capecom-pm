package tickethandler

import (
	ticketsvc "capecom-pm/internal/services/ticket"
	"capecom-pm/internal/utils/bind"
	"capecom-pm/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	HistoryService *ticketsvc.HistoryService
}

func NewHistoryHandler(service *ticketsvc.HistoryService) *HistoryHandler {
	return &HistoryHandler{HistoryService: service}
}

func (h *HistoryHandler) GetTicketHistory(c *gin.Context) {
	ticketUUID := c.Param("ticketId")

	pg, _ := bind.PaginationBinder(c, "get_ticket_history")

	result, err := h.HistoryService.GetAllByTicket(ticketUUID, pg)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}
