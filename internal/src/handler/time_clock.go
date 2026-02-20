package handler

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils"
	"github.com/akashefrath/capecom-pm/internal/utils/bind"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type TimeClockHandler struct {
	TimeClock *service.TimeClock
}

func NewTimeClock(timeClock *service.TimeClock) TimeClockHandler {
	return TimeClockHandler{TimeClock: timeClock}
}

func (h *TimeClockHandler) GetTodayDetails(c *gin.Context) {
	userID := utils.GetUserID(c)

	data, err := h.TimeClock.GetTodayDetails(&userID)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})

}

func (h *TimeClockHandler) ClockIn(c *gin.Context) {
	var req dto.TimeClockRequest
	isValid := bind.AndValidate(c, &req, "time_clock")
	if !isValid {
		return
	}

	employeeID := utils.GetUserID(c)

	data, err := h.TimeClock.ClockIn(employeeID, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONCreated(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *TimeClockHandler) ClockOut(c *gin.Context) {

	var req dto.TimeClockRequest
	isValid := bind.AndValidate(c, &req, "time_clock")
	if !isValid {
		return
	}

	employeeID := utils.GetUserID(c)

	data, err := h.TimeClock.ClockOut(employeeID, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONCreated(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *TimeClockHandler) BreakIn(c *gin.Context) {
	var req dto.TimeClockRequest
	isValid := bind.AndValidate(c, &req, "time_clock")
	if !isValid {
		return
	}

	employeeID := utils.GetUserID(c)

	data, err := h.TimeClock.BreakIn(employeeID, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONCreated(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *TimeClockHandler) BreakOut(c *gin.Context) {
	var req dto.TimeClockRequest
	isValid := bind.AndValidate(c, &req, "time_clock")
	if !isValid {
		return
	}

	employeeID := utils.GetUserID(c)

	data, err := h.TimeClock.BreakOut(employeeID, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONCreated(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
