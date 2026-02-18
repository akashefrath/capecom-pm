package handler

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils"
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

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	isValid := bind.AndValidate(c, &req, "login")
	if !isValid {
		return
	}
	token, err := h.Auth.Login(req.Email, req.Password)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,

		Data: token,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {

	var req dto.RefreshTokenRequest

	isValid := bind.AndValidate(c, &req, "login")
	if !isValid {
		return
	}
	token, err := h.Auth.RefreshToken(req.Token)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,

		Data: token,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	_, userID, _ := utils.GetUserData(c)

	data, err := h.Auth.FindUser(userID)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {

	jti := utils.GetJTI(c)

	err := h.Auth.LogoutUserByJTI(jti)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Message: utils.GetMessageWithExtra("func_success", c, "logout"),
	})
}
