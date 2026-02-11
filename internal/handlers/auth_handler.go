package handlers

import (
	authdto "capecom-pm/internal/domain/dto/auth"
	"capecom-pm/internal/services"
	"capecom-pm/internal/utils"
	"net/http"

	"capecom-pm/internal/utils/bind"
	"capecom-pm/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(AuthService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: AuthService,
	}
}

type FieldErrors map[string][]string

func (l *AuthHandler) Login(c *gin.Context) {
	var req authdto.LoginRequest
	validate := bind.AndValidate(c, &req, "")
	if !validate {
		return
	}
	tokenData, err := l.AuthService.Login(req)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Data: tokenData,
	})
}

func (l *AuthHandler) Refresh(c *gin.Context) {

	var req authdto.RefreshTokenRequest
	if validate := bind.AndValidate(c, &req, ""); !validate {
		return
	}

	tokenData, err := l.AuthService.Refresh(req.Token)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Data: tokenData,
	})
}

func (l *AuthHandler) Me(c *gin.Context) {
	userID := utils.GetUserID(c)
	if usr, err := l.AuthService.FindUserByUuid(userID); err != nil {
		response.FromError(c, err)
	} else {
		response.JSON(c, http.StatusOK, response.APIResponse{
			Success: true,
			Data:    usr,
		})
	}

}

func (l *AuthHandler) Logout(c *gin.Context) {
	jti := utils.GetJTI(c)
	err := l.AuthService.LogoutUserByJTI(jti)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, response.APIResponse{
		Success: true,
		Message: utils.GetMessage("logged_out_success", c),
	})

}
