package handlers

import (
	"capecom-pm/internal/domain/dto"
	"capecom-pm/internal/services"
	"capecom-pm/internal/utils/bind"
	"capecom-pm/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(UserService *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: UserService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	validate := bind.AndValidate(c, &req, "create_user")
	if !validate {
		return
	}
	if err := h.UserService.CreateUser(c, req); err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, 200, response.APIResponse{})
}
