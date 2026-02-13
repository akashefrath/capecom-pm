package dto

import "capecom-pm/internal/domain/models"

type CreateProjectAssetRequest struct {
	Title       string  `json:"title" form:"title" binding:"required,min=2,max=150"`
	AssetType   string  `json:"asset_type" form:"asset_type" binding:"required,oneof=text file image link secret"`
	Description *string `json:"description" form:"description" binding:"omitempty"`
	FileUUID    *string `json:"file_id" form:"file_id" binding:"omitempty,uuid"`
	Content     *string `json:"content" form:"content" binding:"omitempty"`
	IsPrivate   *bool   `json:"is_private" form:"is_private" binding:"omitempty"`
}

type UpdateProjectAssetRequest struct {
	Title       *string `json:"title" form:"title" binding:"omitempty,min=2,max=150"`
	AssetType   *string `json:"asset_type" form:"asset_type" binding:"omitempty,oneof=text file image link secret"`
	Description *string `json:"description" form:"description" binding:"omitempty"`
	FileUUID    *string `json:"file_id" form:"file_id" binding:"omitempty,uuid"`
	Content     *string `json:"content" form:"content" binding:"omitempty"`
	IsPrivate   *bool   `json:"is_private" form:"is_private" binding:"omitempty"`
}

type ProjectAssetResponse struct {
	Id          string  `json:"id"`
	ProjectID   string  `json:"project_id"`
	Title       string  `json:"title"`
	AssetType   string  `json:"asset_type"`
	Description *string `json:"description"`
	FileID      *string `json:"file_id"`
	Content     *string `json:"content"`
	IsPrivate   bool    `json:"is_private"`
	Status      string  `json:"status"`

	models.BaseResponse
}
