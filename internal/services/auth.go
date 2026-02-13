package services

import (
	"capecom-pm/internal/domain/dto"
	authdto "capecom-pm/internal/domain/dto/auth"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	"capecom-pm/internal/utils"
	jwtutil "capecom-pm/internal/utils/jwt"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	AuthRepo    *repositories.AuthRepo
	jwt         *jwtutil.Manager
	userRepo    *repositories.UserRepo
	sessionRepo *repositories.SessionRepo
	redisRepo   *cacherepo.RedisRepo
}

func NewAuthService(AuthRepo *repositories.AuthRepo, jwt *jwtutil.Manager, userRepo *repositories.UserRepo, sessionRepo *repositories.SessionRepo, cacheRepo *cacherepo.RedisRepo) *AuthService {
	return &AuthService{
		AuthRepo:    AuthRepo,
		jwt:         jwt,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		redisRepo:   cacheRepo,
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
		return s.CreateAndReturnToken(usr.UUID, "", usr.IsAdmin)
	}

}

func (s AuthService) Refresh(token string) (*authdto.LoginResponse, error) {
	//userID := s.userRepo.GetActiveUserIDByUuid(userUuid)
	//if userID == nil {
	//	return nil, domainerrors.ErrUserNotFound
	//}

	return s.CreateAndReturnToken("", token, false)

}

func (s AuthService) CreateAndReturnToken(userUuid, oldToken string, isAdmin bool) (*authdto.LoginResponse, error) {

	jti := uuid.NewString()

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	hashedToken := utils.HashToken(refreshToken)

	var accessToken string

	err = s.sessionRepo.DB.Transaction(func(tx *gorm.DB) error {

		if oldToken == "" {
			var userID int64 = 0

			userIdIn, err := s.redisRepo.GetUserIdByUuid(userUuid, *s.userRepo)
			if err != nil || userIdIn == nil {
				return domainerrors.ErrUnauthorized
			}

			userID = *userIdIn

			if userID == int64(0) {
				return domainerrors.ErrUnauthorized
			}

			session := &models.Session{
				UserID:           userID,
				JTI:              jti,
				RefreshTokenHash: hashedToken,
				RefreshExpiresAt: s.jwt.GetExpireTime(),
			}

			if err := s.sessionRepo.Create(session); err != nil {
				return err
			}

		} else {

			token, err := s.sessionRepo.GetByHashedToken(utils.HashToken(oldToken))
			if err != nil || token == nil || token.Status != models.SessionStatusActive {
				return domainerrors.ErrUnauthorized
			}

			// expiry check
			if token.RefreshExpiresAt.Before(time.Now()) {
				return domainerrors.ErrUnauthorized
			}

			now := time.Now()
			token.RefreshTokenHash = hashedToken
			token.LastUsedAt = now
			token.RotatedAt = &now
			token.JTI = jti
			token.RefreshExpiresAt = s.jwt.GetExpireTime()

			if err := s.sessionRepo.Update(token); err != nil {
				return err
			}

			userUuidNew, _ := s.redisRepo.GetUserUuidById(token.UserID, *s.userRepo)
			if userUuidNew == nil {
				return domainerrors.ErrUnauthorized
			}
			userUuid = *userUuidNew

		}

		isAdminIn, err := s.userRepo.IsAdmin(userUuid)

		if err != nil || userUuid == "" {
			return domainerrors.ErrUnauthorized
		}
		tokenType := jwtutil.TokenTypeUser
		if isAdminIn {
			tokenType = jwtutil.TokenTypeAdmin
		}

		accessToken, err = s.jwt.CreateToken(userUuid, tokenType, jti)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &authdto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		IsAdmin:      isAdmin,
	}, nil
}

func (s AuthService) FindUserByUuid(userUuid string) (*dto.UserResponse, error) {
	return s.userRepo.FindByUUID(userUuid)

}

func (s AuthService) LogoutUserByJTI(jti string) error {
	cacheKey := fmt.Sprintf("session:jti:%s", jti)
	_ = s.redisRepo.Delete(context.Background(), cacheKey)

	return s.sessionRepo.RevokeSession(jti)
}
