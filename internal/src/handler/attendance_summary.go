package handler

import (
	"github.com/akashefrath/capecom-pm/internal/domain/common"
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils/bind"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type AttendanceSummaryHandler struct {
	AttendanceSummary *service.AttendanceSummary
}

func NewAttendanceSummary(attendanceSummary *service.AttendanceSummary) AttendanceSummaryHandler {
	return AttendanceSummaryHandler{AttendanceSummary: attendanceSummary}
}

func (h *AttendanceSummaryHandler) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	data, err := h.AttendanceSummary.GetByUUID(uuid)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *AttendanceSummaryHandler) GetList(c *gin.Context) {
	var pg = common.Pagination{}
	_ = bind.QueryBinder(c, &pg, "attendance_summary")
	pg.Normalize()

	var query dto.AttendanceSummaryListQuery
	isValid := bind.AndValidate(c, &query, "attendance_summary")
	if !isValid {
		return
	}

	data, err := h.AttendanceSummary.GetList(pg, query)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
