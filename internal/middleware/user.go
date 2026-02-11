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
	"time"

	"github.com/gin-gonic/gin"
)

type UserMiddleware struct {
	JWTManager  *jwtutil.Manager
	UserRepo    *repositories.UserRepo
	redisRepo   *cacherepo.RedisRepo
	SessionRepo *repositories.SessionRepo
}

func NewUserMiddleware(jwtManager *jwtutil.Manager, userRepo *repositories.UserRepo, redis *cacherepo.RedisRepo, SessionRepo *repositories.SessionRepo) *UserMiddleware {
	return &UserMiddleware{
		JWTManager:  jwtManager,
		UserRepo:    userRepo,
		redisRepo:   redis,
		SessionRepo: SessionRepo,
	}
}

// VerifyUserToken middleware verifies user JWT token and user status
func (m *UserMiddleware) VerifyUserToken() gin.HandlerFunc {
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

		// Validate token with user token type
		claims, err := m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeUser)
		if err != nil {
			response.FromError(c, err)
			c.Abort()
			return
		}

		// Verify user exists and is active
		var status string
		err = m.UserRepo.DB.Model(models.User{}).Where("uuid = ?", claims.UserID).Select("status").Scan(&status).Error
		if err != nil {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		if status != models.SessionStatusActive {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		cacheKey := fmt.Sprintf("session:jti:%s", claims.ID)

		session, err := cacherepo.GetOrSet(
			context.Background(),
			m.redisRepo,
			cacheKey,
			5*time.Minute, // Match access token TTL
			func() (*models.Session, error) {
				return m.SessionRepo.GetByJTI(claims.ID)
			},
		)
		if err != nil || session == nil || session.Status != models.SessionStatusActive {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Set user ID in context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Set("jti", claims.ID)
		c.Next()
	}
}

// VerifyUserRefreshToken middleware verifies user refresh JWT token
func (m *UserMiddleware) VerifyUserRefreshToken() gin.HandlerFunc {
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

		// Validate token with user refresh token type
		claims, err := m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeRefresh)
		if err != nil {
			response.FromError(c, err)
			c.Abort()
			return
		}

		// Verify user exists and is active
		var status string

		err = m.UserRepo.DB.Model(&struct {
			Status string
		}{}).Where("uuid = ?", claims.UserID).Select("status").Scan(&status).Error
		if err != nil {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		if status != models.SessionStatusActive {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}
		cacheKey := fmt.Sprintf("session:jti:%s", claims.ID)
		session, err := cacherepo.GetOrSet(
			context.Background(),
			m.redisRepo,
			cacheKey,
			5*time.Minute, // Match access token TTL
			func() (*models.Session, error) {
				return m.SessionRepo.GetByJTI(claims.ID)
			},
		)
		if err != nil || session == nil || session.Status != models.SessionStatusActive {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}
		// Set user ID in context for downstream handlers
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
