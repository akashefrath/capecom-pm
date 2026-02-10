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
	cacheRepo   *cacherepo.RedisRepo
}

func NewAuthService(AuthRepo *repositories.AuthRepo, jwt *jwtutil.Manager, userRepo *repositories.UserRepo, sessionRepo *repositories.SessionRepo, cacheRepo *cacherepo.RedisRepo) *AuthService {
	return &AuthService{
		AuthRepo:    AuthRepo,
		jwt:         jwt,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		cacheRepo:   cacheRepo,
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

func (s AuthService) CreateAndReturnToken(userUuid, oldToken string) (*authdto.LoginResponse, error) {

	jti := uuid.NewString()

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	hashedToken := utils.HashToken(refreshToken)

	var accessToken string

	err = s.sessionRepo.DB.Transaction(func(tx *gorm.DB) error {

		if oldToken == "" {

			userID, _ := cacherepo.GetCacheDataOrDB(
				func() (*int64, error) {
					data, err := s.cacheRepo.GetString(context.Background(), "id_by_uuid:"+userUuid)
					if err != nil {
						return nil, err
					}
					finalInt, err := utils.ToInt64(data)
					if err != nil {
						return nil, err
					}
					return &finalInt, nil
				},
				func() (*int64, error) {
					return s.userRepo.GetActiveUserIDByUuid(userUuid), nil
				},
				func(data *int64) error {
					return s.cacheRepo.SetString(context.Background(), "id_by_uuid:"+userUuid, data, 0)
				},
			)
			if err != nil {
				return err
			}
			if userID == nil {
				return domainerrors.ErrUnauthorized
			}

			session := &models.Session{
				UserID:           *userID,
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
			token.RefreshExpiresAt = s.jwt.GetExpireTime()

			if err := s.sessionRepo.Update(token); err != nil {
				return err
			}

			userUuidNew, _ := cacherepo.GetCacheDataOrDB(
				func() (*string, error) {
					data, err := s.cacheRepo.GetString(context.Background(), fmt.Sprintf("uuid_by_id:%d", token.UserID))
					if err != nil || data == "" {
						return nil, domainerrors.ErrUnauthorized
					}

					return &data, nil

				},
				func() (*string, error) {

					return s.userRepo.GetActiveUserUuidByID(token.UserID), nil
				},
				func(data *string) error {
					 
					return s.cacheRepo.SetString(context.Background(), fmt.Sprintf("uuid_by_id:%d", token.UserID), data, 0)
				},
			)

			if userUuidNew == nil {
				return domainerrors.ErrUnauthorized
			} else {
				userUuid = *userUuidNew
			}

		}

		var err error
		accessToken, err = s.jwt.CreateToken(userUuid, jwtutil.TokenTypeAdmin, jti)
		return err
	})

	if err != nil {
		return nil, err
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
