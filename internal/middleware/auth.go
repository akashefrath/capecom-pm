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

type AuthMiddleware struct {
	JWTManager  *jwtutil.Manager
	UserRepo    *repositories.UserRepo
	redisRepo   *cacherepo.RedisRepo
	SessionRepo *repositories.SessionRepo
}

func NewAuthMiddleware(jwtManager *jwtutil.Manager, userRepo *repositories.UserRepo, redis *cacherepo.RedisRepo, sessionRepo *repositories.SessionRepo) *AuthMiddleware {
	return &AuthMiddleware{
		JWTManager:  jwtManager,
		UserRepo:    userRepo,
		redisRepo:   redis,
		SessionRepo: sessionRepo,
	}
}

// VerifyToken middleware verifies both user and admin JWT tokens
// Sets userID and isAdmin in context
func (m *AuthMiddleware) VerifyToken() gin.HandlerFunc {
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

		// Try to validate as user token first
		claims, err := m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeUser)
		isAdmin := false

		// If user token validation fails, try admin token
		if err != nil {
			claims, err = m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeAdmin)
			if err != nil {
				response.FromError(c, err)
				c.Abort()
				return
			}
			isAdmin = true
		}

		// Verify user exists and is active
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

		if err != nil {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		if *status != "active" {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Verify session is active
		cacheKey := fmt.Sprintf("session:jti:%s", claims.ID)
		session, err := cacherepo.GetOrSet(
			context.Background(),
			m.redisRepo,
			cacheKey,
			5*time.Minute,
			func() (*models.Session, error) {
				return m.SessionRepo.GetByJTI(claims.ID)
			},
		)
		if err != nil || session == nil || session.Status != models.SessionStatusActive {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Set user ID and admin status in context
		c.Set("userID", claims.UserID)
		c.Set("isAdmin", isAdmin)
		c.Set("jti", claims.ID)

		c.Next()
	}
}
