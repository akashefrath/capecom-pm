package handler

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils/bind"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type ShiftSystemGroupHandler struct {
	ShiftSystemGroup *service.ShiftSystemGroup
}

func NewShiftSystemGroup(shiftSystemGroup *service.ShiftSystemGroup) ShiftSystemGroupHandler {
	return ShiftSystemGroupHandler{ShiftSystemGroup: shiftSystemGroup}
}

func (h *ShiftSystemGroupHandler) Create(c *gin.Context) {
	var req dto.CreateShiftSystemGroupRequest
	isValid := bind.AndValidate(c, &req, "shift_system_group")
	if !isValid {
		return
	}

	group, err := h.ShiftSystemGroup.Create(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONCreated(c, response.APIResponse{
		Success: true,
		Data:    group,
	})
}

func (h *ShiftSystemGroupHandler) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var req dto.UpdateShiftSystemGroupRequest
	isValid := bind.AndValidate(c, &req, "shift_system_group")
	if !isValid {
		return
	}

	group, err := h.ShiftSystemGroup.Update(uuid, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    group,
	})
}

func (h *ShiftSystemGroupHandler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	err := h.ShiftSystemGroup.Delete(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
	})
}

func (h *ShiftSystemGroupHandler) GetAll(c *gin.Context) {
	data, err := h.ShiftSystemGroup.GetAll()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *ShiftSystemGroupHandler) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	data, err := h.ShiftSystemGroup.GetByUUID(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
