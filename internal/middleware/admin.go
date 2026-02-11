package middleware

import (
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	jwtutil "capecom-pm/internal/utils/jwt"
	"capecom-pm/internal/utils/response"
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct {
	JWTManager *jwtutil.Manager
	UserRepo   *repositories.UserRepo
	redisRepo  *cacherepo.RedisRepo
}

func NewAdminMiddleware(jwtManager *jwtutil.Manager, userRepo *repositories.UserRepo, redis *cacherepo.RedisRepo) *AdminMiddleware {
	return &AdminMiddleware{
		JWTManager: jwtManager,
		UserRepo:   userRepo,
		redisRepo:  redis,
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
		err = m.redisRepo.SetString(context.Background(), claims.UserID, status, 0)
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
	userStatusCacheKey := fmt.Sprintf("user_status_%s", claims.UserID)

	status, err := cacherepo.GetOrSet(
		context.Background(),
		m.redisRepo,
		userStatusCacheKey,
		0,
		func() (*string, error) {
			var status string
			err := m.UserRepo.DB.Model(models.User{}).Where("uuid = ?", claims.UserID).Select("status").Scan(&status).Error
			return &status, err
		})
	return status, err
}
