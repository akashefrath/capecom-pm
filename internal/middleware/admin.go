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

type AdminMiddleware struct {
	JWTManager  *jwtutil.Manager
	UserRepo    *repositories.UserRepo
	redisRepo   *cacherepo.RedisRepo
	SessionRepo *repositories.SessionRepo
}

func NewAdminMiddleware(jwtManager *jwtutil.Manager, userRepo *repositories.UserRepo, redis *cacherepo.RedisRepo, SessionRepo *repositories.SessionRepo) *AdminMiddleware {
	return &AdminMiddleware{
		JWTManager:  jwtManager,
		UserRepo:    userRepo,
		redisRepo:   redis,
		SessionRepo: SessionRepo,
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

		if err != nil || status == nil {

			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		if *status != models.SessionStatusActive {

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
			println("aaaaaaa5")
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Set user ID in context for downstream handlers
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

func verifyUserStatus(m *AdminMiddleware, claims *jwtutil.Claims) (*string, error) {
	userStatusCacheKey := fmt.Sprintf("user_status_%s", claims.UserID)

	status, err := cacherepo.GetOrSet(
		context.Background(),
		m.redisRepo,
		userStatusCacheKey,
		0,
		func() (*string, error) {
			var status *string
			status, err := m.UserRepo.FindUserStatus(claims.UserID)
			return status, err
		})
	return status, err
}
