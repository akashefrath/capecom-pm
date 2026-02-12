package handlers

import (
	"capecom-pm/internal/domain/dto"
	"capecom-pm/internal/services"
	"capecom-pm/internal/utils"
	"capecom-pm/internal/utils/bind"
	"capecom-pm/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	ProjectService *services.ProjectService
}

func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{ProjectService: projectService}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if !bind.AndValidate(c, &req, "create_project") {
		return
	}

	userID := utils.GetUserID(c)

	project, err := h.ProjectService.Create(req, userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, response.APIResponse{
		Success: true,
		Data:    project,
	})
}
