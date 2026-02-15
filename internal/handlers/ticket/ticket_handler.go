package tickethandler

import (
	"capecom-pm/internal/domain/dto"
	ticketsvc "capecom-pm/internal/services/ticket"
	"capecom-pm/internal/utils"
	"capecom-pm/internal/utils/bind"
	"capecom-pm/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	TicketService *ticketsvc.TicketService
}

func NewTicketHandler(service *ticketsvc.TicketService) *TicketHandler {
	return &TicketHandler{TicketService: service}
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	projectUUID := c.Param("projectId")

	var req dto.CreateTicketRequest
	if !bind.AndValidate(c, &req, "create_ticket") {
		return
	}

	result, err := h.TicketService.Create(projectUUID, req, utils.GetUserID(c))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, response.APIResponse{Success: true, Data: result})
}

func (h *TicketHandler) GetTickets(c *gin.Context) {
	projectUUID := c.Param("projectId")
	userID := utils.GetUserID(c)

	pg, _ := bind.PaginationBinder(c, "get_tickets")

	result, err := h.TicketService.GetAllByProject(projectUUID, pg, userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TicketHandler) GetTicket(c *gin.Context) {
	ticketUUID := c.Param("ticketId")

	result, err := h.TicketService.GetByUUID(ticketUUID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	ticketUUID := c.Param("ticketId")

	var req dto.UpdateTicketRequest
	if !bind.AndValidate(c, &req, "update_ticket") {
		return
	}

	result, err := h.TicketService.Update(ticketUUID, req)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TicketHandler) UpdateLifecycleStatus(c *gin.Context) {
	ticketUUID := c.Param("ticketId")
	userID := utils.GetUserID(c)

	var req dto.UpdateTicketLifecycleRequest
	if !bind.AndValidate(c, &req, "update_ticket_lifecycle") {
		return
	}

	result, err := h.TicketService.UpdateLifecycleStatus(ticketUUID, req.LifecycleStatus, userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TicketHandler) UpdateAssignee(c *gin.Context) {
	projectUUID := c.Param("projectId")
	ticketUUID := c.Param("ticketId")

	var req dto.UpdateTicketAssigneeRequest
	if !bind.AndValidate(c, &req, "update_ticket_assignee") {
		return
	}

	result, err := h.TicketService.UpdateAssignee(ticketUUID, projectUUID, req.AssignedToUUID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	ticketUUID := c.Param("ticketId")

	err := h.TicketService.Delete(ticketUUID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Message: "ticket_deleted"})
}
