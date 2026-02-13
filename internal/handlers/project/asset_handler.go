package projecthandler

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	projectsvc "capecom-pm/internal/services/project"
	"capecom-pm/internal/utils"
	"capecom-pm/internal/utils/bind"
	"capecom-pm/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	AssetService *projectsvc.AssetService
}

func NewAssetHandler(service *projectsvc.AssetService) *AssetHandler {
	return &AssetHandler{AssetService: service}
}

func (h *AssetHandler) CreateAsset(c *gin.Context) {
	projectUUID := c.Param("projectId")

	var req dto.CreateProjectAssetRequest
	if !bind.AndValidate(c, &req, "create_project_asset") {
		return
	}

	userID := utils.GetUserID(c)

	asset, err := h.AssetService.Create(projectUUID, req, userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, response.APIResponse{
		Success: true,
		Data:    asset,
	})
}

func (h *AssetHandler) GetAssets(c *gin.Context) {
	projectUUID := c.Param("projectId")

	var pg common.Pagination
	_ = bind.QueryBinder(c, &pg, "get_project_assets")
	pg.Normalize()

	assets, err := h.AssetService.GetByProjectUUID(projectUUID, pg)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    assets,
	})
}

func (h *AssetHandler) GetAsset(c *gin.Context) {
	assetUUID := c.Param("assetId")

	asset, err := h.AssetService.GetByUUID(assetUUID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    asset,
	})
}

func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	assetUUID := c.Param("assetId")

	var req dto.UpdateProjectAssetRequest
	if !bind.AndValidate(c, &req, "update_project_asset") {
		return
	}

	asset, err := h.AssetService.Update(assetUUID, req)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    asset,
	})
}

func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	assetUUID := c.Param("assetId")

	err := h.AssetService.Delete(assetUUID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Message: "asset_deleted",
	})
}
