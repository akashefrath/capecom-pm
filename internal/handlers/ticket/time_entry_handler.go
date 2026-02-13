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

type TimeEntryHandler struct {
	TimeEntryService *ticketsvc.TimeEntryService
}

func NewTimeEntryHandler(service *ticketsvc.TimeEntryService) *TimeEntryHandler {
	return &TimeEntryHandler{TimeEntryService: service}
}

func (h *TimeEntryHandler) CreateTimeEntry(c *gin.Context) {
	ticketUUID := c.Param("ticketId")

	var req dto.CreateTimeEntryRequest
	if !bind.AndValidate(c, &req, "create_time_entry") {
		return
	}

	result, err := h.TimeEntryService.Create(ticketUUID, req, utils.GetUserID(c))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, response.APIResponse{Success: true, Data: result})
}

func (h *TimeEntryHandler) GetTimeEntries(c *gin.Context) {
	ticketUUID := c.Param("ticketId")

	pg, _ := bind.PaginationBinder(c, "get_time_entries")

	result, err := h.TimeEntryService.GetAllByTicket(ticketUUID, pg)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TimeEntryHandler) GetTimeEntry(c *gin.Context) {
	entryUUID := c.Param("entryId")

	result, err := h.TimeEntryService.GetByUUID(entryUUID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TimeEntryHandler) UpdateTimeEntry(c *gin.Context) {
	entryUUID := c.Param("entryId")

	var req dto.UpdateTimeEntryRequest
	if !bind.AndValidate(c, &req, "update_time_entry") {
		return
	}

	result, err := h.TimeEntryService.Update(entryUUID, req, utils.GetUserID(c))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TimeEntryHandler) DeleteTimeEntry(c *gin.Context) {
	entryUUID := c.Param("entryId")

	err := h.TimeEntryService.Delete(entryUUID, utils.GetUserID(c))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Message: "time_entry_deleted"})
}
