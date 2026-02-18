package handler

import (
	"net/http"

	"github.com/akashefrath/capecom-pm/internal/domain/common"
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils"
	"github.com/akashefrath/capecom-pm/internal/utils/bind"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	User *service.User
}

func NewUser(user *service.User) UserHandler {
	return UserHandler{User: user}
}

func (h UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest

	if valid := bind.AndValidate(c, &req, "create_user"); !valid {
		return
	}
	_, id, _ := utils.GetUserData(c)

	usr, err := h.User.Create(id, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, response.APIResponse{
		Success: true,
		Data:    usr,
	})

}

func (h UserHandler) GetUsers(c *gin.Context) {
	var pg = common.Pagination{}
	_ = bind.QueryBinder(c, &pg, "get_users")
	pg.Normalize()
	filter, _ := bind.FilterBinder(c)
	users, err := h.User.GetAll(pg, filter)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    users,
	})
}

func (h UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.User.GetByID(id)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    user,
	})
}
