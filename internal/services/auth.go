package services

import (
	authdto "capecom-pm/internal/domain/dto/auth"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/repositories"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	AuthRepo *repositories.AuthRepo
}

func NewAuthService(AuthRepo *repositories.AuthRepo) *AuthService {
	return &AuthService{
		AuthRepo: AuthRepo,
	}
}

func (s AuthService) Login(c *gin.Context, req authdto.LoginRequest) error {
	if usr, err := s.AuthRepo.FindUserByEmailAndReturnPassword(req.Email); err != nil {
		return nil
	} else if usr == nil {
		return domainerrors.ErrInvalidCredentials
	} else {

	}

	return nil
}
