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

type FileHandler struct {
	FileService *services.FileService
}

func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{FileService: fileService}
}

func (h *FileHandler) CreateFile(c *gin.Context) {
	var req dto.CreateFileRequest
	if !bind.AndValidate(c, &req, "file") {
		return
	}

	userID := utils.GetUserID(c)

	result, err := h.FileService.CreateFileAndGetUploadURL(userID, req)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, response.APIResponse{
		Success: true,
		Data:    result,
	})
}

func (h *FileHandler) ConfirmUpload(c *gin.Context) {
	var req dto.ConfirmUploadRequest
	if !bind.AndValidate(c, &req, "file") {
		return
	}

	results, err := h.FileService.ConfirmUpload(req.FileIDs)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    results,
	})
}
