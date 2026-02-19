package handler

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils/bind"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type AttendancePolicyHandler struct {
	AttendancePolicy *service.AttendancePolicy
}

func NewAttendancePolicy(attendancePolicy *service.AttendancePolicy) AttendancePolicyHandler {
	return AttendancePolicyHandler{AttendancePolicy: attendancePolicy}
}

func (h *AttendancePolicyHandler) Create(c *gin.Context) {
	var req dto.CreateAttendancePolicyRequest
	isValid := bind.AndValidate(c, &req, "attendance_policy")
	if !isValid {
		return
	}

	policy, err := h.AttendancePolicy.Create(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONCreated(c, response.APIResponse{
		Success: true,
		Data:    policy,
	})
}

func (h *AttendancePolicyHandler) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var req dto.CreateAttendancePolicyRequest
	isValid := bind.AndValidate(c, &req, "attendance_policy")
	if !isValid {
		return
	}

	policy, err := h.AttendancePolicy.Update(uuid, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    policy,
	})
}

func (h *AttendancePolicyHandler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	err := h.AttendancePolicy.Delete(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
	})
}

func (h *AttendancePolicyHandler) GetAll(c *gin.Context) {

	data, err := h.AttendancePolicy.GetAll()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
func (h *AttendancePolicyHandler) GetAllUtils(c *gin.Context) {

	data, err := h.AttendancePolicy.GetAllUtils()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *AttendancePolicyHandler) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	data, err := h.AttendancePolicy.GetByUUID(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *AttendancePolicyHandler) SetDefault(c *gin.Context) {
	uuid := c.Param("uuid")
	data, err := h.AttendancePolicy.SetDefault(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
