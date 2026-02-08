package middleware

import (
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/repositories"
	jwtutil "capecom-pm/internal/utils/jwt"
	"capecom-pm/internal/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct {
	JWTManager *jwtutil.JWTManager
	UserRepo   *repositories.UserRepo
}

func NewAdminMiddleware(jwtManager *jwtutil.JWTManager, userRepo *repositories.UserRepo) *AdminMiddleware {
	return &AdminMiddleware{
		JWTManager: jwtManager,
		UserRepo:   userRepo,
	}
}

// VerifyAdminToken middleware verifies admin JWT token
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

		// Set user ID in context for downstream handlers
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
