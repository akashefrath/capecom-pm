package services

import (
	authdto "capecom-pm/internal/domain/dto/auth"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/repositories"
	"capecom-pm/internal/utils"
	jwtutil "capecom-pm/internal/utils/jwt"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	AuthRepo *repositories.AuthRepo
	jwt      *jwtutil.JWTManager
}

func NewAuthService(AuthRepo *repositories.AuthRepo, jwt *jwtutil.JWTManager) *AuthService {
	return &AuthService{
		AuthRepo: AuthRepo,
		jwt:      jwt,
	}
}

func (s AuthService) Login(c *gin.Context, req authdto.LoginRequest) (*authdto.LoginResponse, error) {

	if usr, err := s.AuthRepo.FindUserByEmailAndReturnPassword(req.Email); err != nil {
		return nil, err
	} else if usr == nil {
		return nil, domainerrors.ErrInvalidCredentials
	} else {

		if !utils.CheckPassword(usr.PasswordHash, req.Password) {
			return nil, domainerrors.ErrInvalidCredentials
		}

		accessToken, _ := s.jwt.CreateToken(usr.UUID, jwtutil.TokenTypeUser)
		refreshToken, _ := s.jwt.CreateToken(usr.UUID, jwtutil.TokenTypeRefresh)
		return &authdto.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
		}, nil
	}

}
