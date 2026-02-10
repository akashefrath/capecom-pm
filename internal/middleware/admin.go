package middleware

import (
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	jwtutil "capecom-pm/internal/utils/jwt"
	"capecom-pm/internal/utils/response"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct {
	JWTManager *jwtutil.Manager
	UserRepo   *repositories.UserRepo
	Redis      *cacherepo.RedisRepo
}

func NewAdminMiddleware(jwtManager *jwtutil.Manager, userRepo *repositories.UserRepo, redis *cacherepo.RedisRepo) *AdminMiddleware {
	return &AdminMiddleware{
		JWTManager: jwtManager,
		UserRepo:   userRepo,
		Redis:      redis,
	}
}

// VerifyAdminToken middleware verifies admin JWT token and user status
func (m *AdminMiddleware) VerifyAdminToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.FromError(c, domainerrors.ErrInvalidToken)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token with admin token type
		claims, err := m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeAdmin)
		if err != nil {
			response.FromError(c, err)
			c.Abort()
			return
		}

		status, err := verifyUserStatus(m, claims)

		if err != nil {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		if status != "active" {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Set user ID in context for downstream handlers
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

// VerifyAdminRefreshToken middleware verifies admin refresh JWT token
func (m *AdminMiddleware) VerifyAdminRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.FromError(c, domainerrors.ErrInvalidToken)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token with admin refresh token type
		claims, err := m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeAdminRefresh)
		if err != nil {
			response.FromError(c, err)
			c.Abort()
			return
		}

		// Verify user exists and is active
		status, _ := verifyUserStatus(m, claims)

		if status != "active" {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}
		err = m.Redis.SetString(context.Background(), claims.UserID, status, 0)
		if err != nil {
			print(err)
			return
		}
		// Set user ID in context for downstream handlers
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

func verifyUserStatus(m *AdminMiddleware, claims *jwtutil.Claims) (any, error) {
	status, err := cacherepo.GetCacheDataOrDB(
		func() (*string, error) {

			data, err := m.Redis.GetString(context.Background(), claims.UserID)
			if err != nil {
				return nil, err
			}

			return &data, nil
		},
		func() (*string, error) {
			var status string
			err := m.UserRepo.DB.Model(models.User{}).Where("uuid = ?", claims.UserID).Select("status").Scan(&status).Error
			return &status, err
		},
		func(v *string) error {
			go func() {
				_ = m.Redis.SetString(context.Background(), claims.UserID, v, 0)
			}()
			return nil
		},
	)
	return status, err
}
