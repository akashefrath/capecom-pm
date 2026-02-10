package services

import (
	"capecom-pm/internal/domain/dto"
	authdto "capecom-pm/internal/domain/dto/auth"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	"capecom-pm/internal/utils"
	jwtutil "capecom-pm/internal/utils/jwt"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	AuthRepo    *repositories.AuthRepo
	jwt         *jwtutil.Manager
	userRepo    *repositories.UserRepo
	sessionRepo *repositories.SessionRepo
}

func NewAuthService(AuthRepo *repositories.AuthRepo, jwt *jwtutil.Manager, userRepo *repositories.UserRepo, sessionRepo *repositories.SessionRepo) *AuthService {
	return &AuthService{
		AuthRepo:    AuthRepo,
		jwt:         jwt,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
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
		return s.CreateAndReturnToken(usr.UUID, "")
	}

}

func (s AuthService) Refresh(token string) (*authdto.LoginResponse, error) {
	//userID := s.userRepo.GetActiveUserIDByUuid(userUuid)
	//if userID == nil {
	//	return nil, domainerrors.ErrUserNotFound
	//}

	return s.CreateAndReturnToken("", token)

}

func (s AuthService) CreateAndReturnToken(userUuid string, oldToken string) (*authdto.LoginResponse, error) {
	jti := uuid.NewString()

	refreshToken, _ := utils.GenerateRefreshToken()
	hashedToken := utils.HashToken(refreshToken)
	var accessToken string

	if oldToken == "" {
		userData := s.userRepo.GetActiveUserIDByUuid(userUuid)
		println(*userData)
		session := &models.Session{
			UserID:           *userData,
			JTI:              jti,
			RefreshTokenHash: hashedToken,
			RefreshExpiresAt: s.jwt.GetExpireTime(),
		}

		err := s.sessionRepo.Create(session)
		if err != nil {
			return nil, err
		}
		accessToken, _ = s.jwt.CreateToken(userUuid, jwtutil.TokenTypeAdmin, jti)
	} else {
		token, err := s.sessionRepo.GetByHashedToken(utils.HashToken(oldToken))
		if err != nil {
			return nil, err
		}
		userData := s.userRepo.GetActiveUserUuidByID(token.UserID)
		token.RefreshTokenHash = hashedToken
		token.LastUsedAt = time.Now()

		err = s.sessionRepo.Update(token)
		if err != nil {
			return nil, err
		}
		accessToken, _ = s.jwt.CreateToken(*userData, jwtutil.TokenTypeAdmin, jti)

	}

	return &authdto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil

}

func (s AuthService) FindUserByUuid(userUuid string) (*dto.UserResponse, error) {
	return s.userRepo.FindByUUID(userUuid)

}
