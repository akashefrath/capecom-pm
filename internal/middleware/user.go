package middleware

import (
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	jwtutil "capecom-pm/internal/utils/jwt"
	"capecom-pm/internal/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserMiddleware struct {
	JWTManager *jwtutil.Manager
	UserRepo   *repositories.UserRepo
	Redis      *cacherepo.RedisRepo
}

func NewUserMiddleware(jwtManager *jwtutil.Manager, userRepo *repositories.UserRepo, redis *cacherepo.RedisRepo) *UserMiddleware {
	return &UserMiddleware{
		JWTManager: jwtManager,
		UserRepo:   userRepo,
		Redis:      redis,
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
