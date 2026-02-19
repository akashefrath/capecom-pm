package handler

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils/bind"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type ShiftSystemHandler struct {
	ShiftSystem *service.ShiftSystem
}

func NewShiftSystem(shiftSystem *service.ShiftSystem) ShiftSystemHandler {
	return ShiftSystemHandler{ShiftSystem: shiftSystem}
}

func (h *ShiftSystemHandler) Create(c *gin.Context) {
	var req dto.CreateShiftSystemRequest
	isValid := bind.AndValidate(c, &req, "shift_system")
	if !isValid {
		return
	}

	shift, err := h.ShiftSystem.Create(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONCreated(c, response.APIResponse{
		Success: true,
		Data:    shift,
	})
}

func (h *ShiftSystemHandler) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var req dto.UpdateShiftSystemRequest
	isValid := bind.AndValidate(c, &req, "shift_system")
	if !isValid {
		return
	}

	shift, err := h.ShiftSystem.Update(uuid, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    shift,
	})
}

func (h *ShiftSystemHandler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	err := h.ShiftSystem.Delete(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
	})
}

func (h *ShiftSystemHandler) GetAll(c *gin.Context) {
	data, err := h.ShiftSystem.GetAll()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *ShiftSystemHandler) GetAllUtils(c *gin.Context) {
	data, err := h.ShiftSystem.GetAllUtils()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *ShiftSystemHandler) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	data, err := h.ShiftSystem.GetByUUID(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *ShiftSystemHandler) SetDefault(c *gin.Context) {
	uuid := c.Param("uuid")
	data, err := h.ShiftSystem.SetDefault(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
