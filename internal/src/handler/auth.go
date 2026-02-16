package handler

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils/bind"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Auth *service.Auth
}

func NewAuth(auth *service.Auth) AuthHandler {
	return AuthHandler{
		Auth: auth,
	}
}

func (a *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	isValid := bind.AndValidate(c, &req, "login")
	if !isValid {
		return
	}
	token, err := a.Auth.Login(req.Email, req.Password)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,

		Data: token,
	})
}
