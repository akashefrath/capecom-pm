package handlers

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	"capecom-pm/internal/services"
	"capecom-pm/internal/utils"
	"capecom-pm/internal/utils/bind"
	"capecom-pm/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	ClientService *services.ClientService
}

func NewClientHandler(clientService *services.ClientService) *ClientHandler {
	return &ClientHandler{ClientService: clientService}
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	var req dto.CreateClientRequest
	if !bind.AndValidate(c, &req, "create_client") {
		return
	}

	userID := utils.GetUserID(c)

	client, err := h.ClientService.Create(req, userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, response.APIResponse{
		Success: true,
		Data:    client,
	})
}

func (h *ClientHandler) GetClientByID(c *gin.Context) {
	uuid := c.Param("id")

	client, err := h.ClientService.GetByUUID(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    client,
	})
}

func (h *ClientHandler) GetClients(c *gin.Context) {
	var pg common.Pagination
	_ = bind.QueryBinder(c, &pg, "get_clients")
	pg.Normalize()

	clients, err := h.ClientService.GetClients(pg)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    clients,
	})
}

func (h *ClientHandler) UpdateClient(c *gin.Context) {
	uuid := c.Param("id")

	var req dto.UpdateClientRequest
	if !bind.AndValidate(c, &req, "update_client") {
		return
	}

	client, err := h.ClientService.Update(uuid, req)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    client,
	})
}
