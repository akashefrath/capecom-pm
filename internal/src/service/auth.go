package service

import (
	"errors"
	"net/http"

	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	domainerrors "github.com/akashefrath/capecom-pm/internal/domain/error"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
	"github.com/akashefrath/capecom-pm/internal/utils"
	jwtutil "github.com/akashefrath/capecom-pm/internal/utils/jwt"
)

type Auth struct {
	Session    *repository.Session
	User       *repository.User
	JWTManager *jwtutil.Manager
}

func NewAuth(auth *repository.Session, user *repository.User, jwtManager *jwtutil.Manager) *Auth {
	return &Auth{
		Session:    auth,
		User:       user,
		JWTManager: jwtManager,
	}
}

func (s *Auth) Login(email, password string) (*dto.LoginResponse, error) {

	errorInvalid := domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrInvalidCredentials.Error(), "Login", "login_error")
	user, err := s.User.GetUserMinimalByEmail(email)
	if err != nil {
		if errors.Is(err, domainerrors.ErrUserNotFound) {
			return nil, errorInvalid
		}
		return nil, domainerrors.Internal("Login", "login_error")
	}

	isValid := utils.CheckPassword(user.Password, password)
	if !isValid {
		return nil, errorInvalid
	}

	return s.CreateAndReturnSession(user.UUID, "")
}

func (s *Auth) CreateAndReturnSession(userUuid string, jti string) (*dto.LoginResponse, error) {

	token, expiresIn, err := s.JWTManager.CreateToken(userUuid, jwtutil.TokenTypeUser, jti)
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	hashedToken := utils.HashToken(refreshToken)
	_, err = s.Session.CreateSession(userUuid, jti, hashedToken)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{
		AccessToken:  token,
		RefreshToken: hashedToken,
		TokenType:    "Bearer",
		IsAdmin:      false,
		ExpiresIn:    expiresIn,
	}, err
}
