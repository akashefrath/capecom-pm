package handlers

import (
	"capecom-pm/internal/services"
	"capecom-pm/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UtilsHandler struct {
	UtilsService *services.UtilsService
}

func NewUtilsHandler(utilsService *services.UtilsService) *UtilsHandler {
	return &UtilsHandler{UtilsService: utilsService}
}

func (h *UtilsHandler) GetAll(c *gin.Context) {
	data, err := h.UtilsService.GetAll()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: data})
}

func (h *UtilsHandler) GetRoles(c *gin.Context) {
	data, err := h.UtilsService.GetRoles()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: data})
}

func (h *UtilsHandler) GetUserGroups(c *gin.Context) {
	data, err := h.UtilsService.GetUserGroups()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: data})
}

func (h *UtilsHandler) GetDesignations(c *gin.Context) {
	data, err := h.UtilsService.GetDesignations()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: data})
}

func (h *UtilsHandler) GetDepartments(c *gin.Context) {
	data, err := h.UtilsService.GetDepartments()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: data})
}

func (h *UtilsHandler) GetClients(c *gin.Context) {
	data, err := h.UtilsService.GetClients()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: data})
}

func (h *UtilsHandler) GetTicketTypes(c *gin.Context) {
	data, err := h.UtilsService.GetTicketTypes()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: data})
}

func (h *UtilsHandler) GetUsers(c *gin.Context) {
	data, err := h.UtilsService.GetUsers()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: data})
}
