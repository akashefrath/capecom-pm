package handlers

import (
	"capecom-pm/internal/domain/dto"
	"capecom-pm/internal/services"
	"capecom-pm/internal/utils/bind"
	"capecom-pm/internal/utils/response"
	"net/http"

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
	if usr, err := h.UserService.CreateUser(req); err != nil {
		response.FromError(c, err)
		return
	} else {
		response.JSON(c, http.StatusCreated, response.APIResponse{
			Success: true,
			Data:    usr,
		})
	}

}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    user,
	})
}
