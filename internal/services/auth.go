package services

import (
	"capecom-pm/internal/domain/dto"
	authdto "capecom-pm/internal/domain/dto/auth"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/repositories"
	"capecom-pm/internal/utils"
	jwtutil "capecom-pm/internal/utils/jwt"

	"github.com/google/uuid"
)

type AuthService struct {
	AuthRepo *repositories.AuthRepo
	jwt      *jwtutil.Manager
	userRepo *repositories.UserRepo
}

func NewAuthService(AuthRepo *repositories.AuthRepo, jwt *jwtutil.Manager, userRepo *repositories.UserRepo) *AuthService {
	return &AuthService{
		AuthRepo: AuthRepo,
		jwt:      jwt,
		userRepo: userRepo,
	}
}

func (s AuthService) Login(req authdto.LoginRequest) (*authdto.LoginResponse, error) {

	if usr, err := s.userRepo.FindByEmail(req.Email); err != nil {
		return nil, err
	} else if usr == nil {
		return nil, domainerrors.ErrInvalidCredentials
	} else {

		if !utils.CheckPassword(usr.PasswordHash, req.Password) {
			return nil, domainerrors.ErrInvalidCredentials
		}
		return s.CreateAndReturnToken(usr.UUID)
	}

}

func (s AuthService) Refresh(userUuid string) (*authdto.LoginResponse, error) {
	userID := s.userRepo.GetActiveUserIDByUuid(userUuid)
	if userID == nil {
		return nil, domainerrors.ErrUserNotFound
	}

	return s.CreateAndReturnToken(userUuid)

}

func (s AuthService) CreateAndReturnToken(userUuid string) (*authdto.LoginResponse, error) {
	jti := uuid.NewString()
	accessToken, _ := s.jwt.CreateToken(userUuid, jwtutil.TokenTypeAdmin, jti)
	refreshToken, _ := s.jwt.CreateToken(userUuid, jwtutil.TokenTypeAdminRefresh, jti)

	refreshToken2, _ := utils.GenerateRefreshToken()
	hashedToken := utils.HashToken(refreshToken2)

	println(hashedToken)
	return &authdto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil

}

func (s AuthService) FindUserByUuid(userUuid string) (*dto.UserResponse, error) {
	return s.userRepo.FindByUUID(userUuid)

}
