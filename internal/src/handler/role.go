package handler

import (
	utilsservice "github.com/akashefrath/capecom-pm/internal/src/service/utils"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	Role *utilsservice.Role
}

func NewRole(role *utilsservice.Role) RoleHandler {
	return RoleHandler{Role: role}
}

func (h *RoleHandler) GetAllActive(c *gin.Context) {
	data, err := h.Role.GetAllActive()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
