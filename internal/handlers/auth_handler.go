package handlers

import (
	authdto "capecom-pm/internal/domain/dto/auth"
	"capecom-pm/internal/services"
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
	err := l.AuthService.Login(c, req)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.JSON(c, 200, response.APIResponse{})
}
