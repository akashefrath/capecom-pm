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

type ProjectHandler struct {
	ProjectService *projectsvc.ProjectService
}

func NewProjectHandler(projectService *projectsvc.ProjectService) *ProjectHandler {
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

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	uuid := c.Param("id")

	var req dto.UpdateProjectRequest
	if !bind.AndValidate(c, &req, "update_project") {
		return
	}

	project, err := h.ProjectService.Update(uuid, req)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    project,
	})
}

func (h *ProjectHandler) UpdateLifecycleStatus(c *gin.Context) {
	uuid := c.Param("id")

	var req dto.UpdateProjectLifecycleRequest
	if !bind.AndValidate(c, &req, "update_project_lifecycle") {
		return
	}

	project, err := h.ProjectService.UpdateLifecycleStatus(uuid, req.LifecycleStatus)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    project,
	})
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	uuid := c.Param("id")

	err := h.ProjectService.Delete(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Message: "project_deleted",
	})
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {

	pg, err := bind.PaginationBinder(c, "get_projects")
	if err != nil {
		/// no need
	}
	userID := utils.GetUserID(c)

	data, err := h.ProjectService.GetProjects(pg, userID)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectId := c.Param("projectId")
	userID := utils.GetUserID(c)

	data, err := h.ProjectService.GetProject(projectId, userID)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
