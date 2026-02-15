package projecthandler

import (
	"capecom-pm/internal/domain/dto"
	projectsvc "capecom-pm/internal/services/project"
	"capecom-pm/internal/utils"
	"capecom-pm/internal/utils/bind"
	"capecom-pm/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	TeamService *projectsvc.TeamService
}

func NewTeamHandler(service *projectsvc.TeamService) *TeamHandler {
	return &TeamHandler{TeamService: service}
}

// --- Managers ---

func (h *TeamHandler) AssignManagers(c *gin.Context) {
	projectUUID := c.Param("projectId")

	var req dto.AssignManagersRequest
	if !bind.AndValidate(c, &req, "assign_managers") {
		return
	}

	result, err := h.TeamService.AssignManagers(projectUUID, req, utils.GetUserID(c))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TeamHandler) RemoveManagers(c *gin.Context) {
	projectUUID := c.Param("projectId")

	var req dto.AssignManagersRequest
	if !bind.AndValidate(c, &req, "remove_managers") {
		return
	}

	result, err := h.TeamService.RemoveManagers(projectUUID, req)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TeamHandler) GetManagers(c *gin.Context) {
	projectUUID := c.Param("projectId")

	result, err := h.TeamService.GetManagers(projectUUID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

// --- Members ---

func (h *TeamHandler) AssignMembers(c *gin.Context) {
	projectUUID := c.Param("projectId")

	var req dto.AssignMembersRequest
	if !bind.AndValidate(c, &req, "assign_members") {
		return
	}

	result, err := h.TeamService.AssignMembers(projectUUID, req, utils.GetUserID(c))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TeamHandler) RemoveMembers(c *gin.Context) {
	projectUUID := c.Param("projectId")

	var req dto.RemoveMembersRequest
	if !bind.AndValidate(c, &req, "remove_members") {
		return
	}

	result, err := h.TeamService.RemoveMembers(projectUUID, req)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}

func (h *TeamHandler) GetMembers(c *gin.Context) {
	projectUUID := c.Param("projectId")

	result, err := h.TeamService.GetMembers(projectUUID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{Success: true, Data: result})
}
