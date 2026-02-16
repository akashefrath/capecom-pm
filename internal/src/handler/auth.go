package handler

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	Auth *service.Auth
}

func NewAuth(auth *service.Auth) AuthHandler {
	return AuthHandler{
		Auth: auth,
	}
}

func (a *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	err := c.Bind().Body(&req)
	if err != nil {
		return err
	}

	err = a.Auth.Login(req.Email, req.Password)
	if err != nil {
		return err
	}
	return response.JSONOk(c, response.APIResponse{
		Success: true,
		Message: "Pongs",
	})
}
