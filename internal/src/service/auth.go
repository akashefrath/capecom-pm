package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	domainerrors "github.com/akashefrath/capecom-pm/internal/domain/error"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
	utilsrepository "github.com/akashefrath/capecom-pm/internal/src/repository/utils"
	"github.com/akashefrath/capecom-pm/internal/utils"
	jwtutil "github.com/akashefrath/capecom-pm/internal/utils/jwt"
	"github.com/google/uuid"
)

type Auth struct {
	Session    *repository.Session
	User       *repository.User
	JWTManager *jwtutil.Manager
	Redis      *utilsrepository.Redis
}

func NewAuth(auth *repository.Session, user *repository.User, jwtManager *jwtutil.Manager, Redis *utilsrepository.Redis) *Auth {
	return &Auth{
		Session:    auth,
		User:       user,
		JWTManager: jwtManager,
		Redis:      Redis,
	}
}

func (s *Auth) Login(email, password string) (*dto.LoginResponse, error) {

	errorInvalid := domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrInvalidCredentials.Error(), "Login", "login_error")
	user, err := s.User.GetUserMinimalByEmail(email)

	if err != nil {
		if errors.Is(err, domainerrors.ErrUserNotFound) {
			return nil, errorInvalid
		}
		return nil, errorInvalid
	}

	isValid := utils.CheckPassword(user.Password, password)
	if !isValid {
		return nil, errorInvalid
	}
	jti := uuid.NewString()

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, errorInvalid
	}
	hashedToken := utils.HashToken(refreshToken)
	_, err = s.Session.CreateSession(user.UUID, jti, hashedToken)
	if err != nil {
		return nil, errorInvalid
	}

	return s.CreateAndReturnSession(user.UUID, jti, refreshToken, user.IsAdmin)
}

func (s *Auth) RefreshToken(oldRefreshToken string) (*dto.LoginResponse, error) {
	hashedOldToken := utils.HashToken(oldRefreshToken)
	_, id, userUuid, isAdmin, err := s.Session.GetSessionJTI(hashedOldToken)
	if err != nil {
		return nil, domainerrors.ErrUnauthorized
	}
	refreshToken, err := utils.GenerateRefreshToken()
	hashedToken := utils.HashToken(refreshToken)
	jti := uuid.NewString()
	err = s.Session.UpdateSession(*id, hashedToken, jti)
	if err != nil {
		return nil, domainerrors.ErrUnauthorized
	}
	return s.CreateAndReturnSession(*userUuid, jti, refreshToken, *isAdmin)
}

func (s *Auth) FindUser(id int64) (*dto.User, error) {
	return s.User.FinUserByID(id)

}

func (s *Auth) LogoutUserByJTI(jti string) error {
	count, _ := s.Session.RevokeTokenByJti(jti)
	zeroCount := int64(0)
	if count == zeroCount {
		return domainerrors.ErrUnauthorized
	}
	cacheKey := fmt.Sprintf("session:jti:%s", jti)
	err := s.Redis.Delete(context.Background(), cacheKey)
	if err != nil {
		return err
	}
	return nil
}

func (s *Auth) CreateAndReturnSession(userUuid string, jti string, refreshToken string, isAdmin bool) (*dto.LoginResponse, error) {
	tokenType := jwtutil.TokenTypeUser
	if isAdmin {
		tokenType = jwtutil.TokenTypeAdmin
	}
	token, expiresIn, err := s.JWTManager.CreateToken(userUuid, tokenType, jti)
	return &dto.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		IsAdmin:      false,
		ExpiresIn:    expiresIn,
	}, err
}
