package handler

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils/bind"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type AttendancePolicyGroupHandler struct {
	AttendancePolicyGroup *service.AttendancePolicyGroup
}

func NewAttendancePolicyGroup(attendancePolicyGroup *service.AttendancePolicyGroup) AttendancePolicyGroupHandler {
	return AttendancePolicyGroupHandler{AttendancePolicyGroup: attendancePolicyGroup}
}

func (h *AttendancePolicyGroupHandler) Create(c *gin.Context) {
	var req dto.CreateAttendancePolicyGroupRequest
	isValid := bind.AndValidate(c, &req, "attendance_policy_group")
	if !isValid {
		return
	}

	group, err := h.AttendancePolicyGroup.Create(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONCreated(c, response.APIResponse{
		Success: true,
		Data:    group,
	})
}

func (h *AttendancePolicyGroupHandler) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var req dto.UpdateAttendancePolicyGroupRequest
	isValid := bind.AndValidate(c, &req, "attendance_policy_group")
	if !isValid {
		return
	}

	group, err := h.AttendancePolicyGroup.Update(uuid, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    group,
	})
}

func (h *AttendancePolicyGroupHandler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	err := h.AttendancePolicyGroup.Delete(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
	})
}

func (h *AttendancePolicyGroupHandler) GetAll(c *gin.Context) {
	data, err := h.AttendancePolicyGroup.GetAll()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *AttendancePolicyGroupHandler) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	data, err := h.AttendancePolicyGroup.GetByUUID(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *AttendancePolicyGroupHandler) AssignUsers(c *gin.Context) {
	uuid := c.Param("uuid")
	var req dto.AssignUsersToGroupRequest
	isValid := bind.AndValidate(c, &req, "attendance_policy_group")
	if !isValid {
		return
	}

	err := h.AttendancePolicyGroup.AssignUsers(uuid, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
	})
}

func (h *AttendancePolicyGroupHandler) RemoveUsers(c *gin.Context) {
	var req dto.RemoveUsersFromGroupRequest
	isValid := bind.AndValidate(c, &req, "attendance_policy_group")
	if !isValid {
		return
	}

	err := h.AttendancePolicyGroup.RemoveUsers(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
	})
}

func (h *AttendancePolicyGroupHandler) GetUsersInGroup(c *gin.Context) {
	uuid := c.Param("uuid")
	data, err := h.AttendancePolicyGroup.GetUsersInGroup(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
