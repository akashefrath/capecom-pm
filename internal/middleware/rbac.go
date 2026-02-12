package middleware

import (
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	"capecom-pm/internal/utils"
	jwtutil "capecom-pm/internal/utils/jwt"
	"capecom-pm/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RABCMiddleware struct {
	JWTManager  *jwtutil.Manager
	UserRepo    *repositories.UserRepo
	redisRepo   *cacherepo.RedisRepo
	SessionRepo *repositories.SessionRepo
}

func NewRABCMiddleware(jwtManager *jwtutil.Manager, userRepo *repositories.UserRepo, redis *cacherepo.RedisRepo, sessionRepo *repositories.SessionRepo) *RABCMiddleware {
	return &RABCMiddleware{
		JWTManager:  jwtManager,
		UserRepo:    userRepo,
		redisRepo:   redis,
		SessionRepo: sessionRepo,
	}
}

func (m *RABCMiddleware) IsManagerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := utils.GetUserID(c)
		is, err := m.UserRepo.IsManagerOrAdmin(userID)
		if err != nil || !is {
			response.FromError(c, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "check_role", "access_1"))
			c.Abort()
		}
		c.Next()
	}

}
