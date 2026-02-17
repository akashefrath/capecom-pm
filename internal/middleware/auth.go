package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	domainerrors "github.com/akashefrath/capecom-pm/internal/domain/error"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
	utilsrepository "github.com/akashefrath/capecom-pm/internal/src/repository/utils"
	"github.com/akashefrath/capecom-pm/internal/utils"
	jwtutil "github.com/akashefrath/capecom-pm/internal/utils/jwt"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type AllowedType string

const (
	AllowedTypeUser  AllowedType = "user"
	AllowedTypeAdmin AllowedType = "admin"
	AllowedTypeAll   AllowedType = "all"
)

type Auth struct {
	JWTManager  *jwtutil.Manager
	UserRepo    *repository.User
	redisRepo   *utilsrepository.Redis
	SessionRepo *repository.Session
}

func NewAuth(jwtManager *jwtutil.Manager, userRepo *repository.User, redis *utilsrepository.Redis, sessionRepo *repository.Session) *Auth {
	return &Auth{
		JWTManager:  jwtManager,
		UserRepo:    userRepo,
		redisRepo:   redis,
		SessionRepo: sessionRepo,
	}
}

// VerifyToken middleware verifies both user and admin JWT tokens
// Sets userID and isAdmin in context
func (m *Auth) VerifyToken(allowedType AllowedType) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.FromError(c, domainerrors.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenString := authHeader[7:]
		var claims *jwtutil.Claims
		var err error
		isAdmin := false

		// 1. Efficiency: Only attempt the validation relevant to the route
		switch allowedType {
		case AllowedTypeAdmin:
			claims, err = m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeAdmin)
			println("ad")
			isAdmin = true
		case AllowedTypeUser:
			claims, err = m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeUser)
			println("user")
		case AllowedTypeAll:
			// For "All", we try User first as it's the most common
			claims, err = m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeUser)
			if err != nil {
				// Fallback to Admin if User validation fails
				claims, err = m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeAdmin)
				isAdmin = err == nil
			}
		}

		if err != nil {
			response.FromError(c, domainerrors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 2. Optimized Session & Status Lookup
		// Using c.Request.Context() ensures the DB/Redis call cancels if the client hangs up
		cacheKey := fmt.Sprintf("session:jti:%s", claims.ID)
		session, err := utilsrepository.GetOrSet(
			c.Request.Context(),
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

		///// verify user status
		//status, err := verifyUserStatus(m.UserRepo, m.redisRepo, &session.UserID)
		//
		//if err != nil || status == nil {
		//	response.FromError(c, err)
		//	c.Abort()
		//	return
		//
		//}
		//
		//if *status != models.SessionStatusActive {
		//	response.FromError(c, domainerrors.ErrUnauthorized)
		//	c.Abort()
		//	return
		//}

		// Set context values
		c.Set(utils.CtxUserUUID, claims.UserID)
		c.Set(utils.CtxUserID, session.UserID)
		c.Set(utils.CtxIsAdmin, isAdmin)
		c.Set(utils.CtxJTI, claims.ID)

		c.Next()
	}
}
func verifyUserStatus(UserRepo *repository.User, redisRepo *utilsrepository.Redis, userID *int64) (*string, error) {
	userStatusCacheKey := fmt.Sprintf("user_status_%d", userID)
	status, err := utilsrepository.GetOrSet(
		context.Background(),
		redisRepo,
		userStatusCacheKey,
		0,
		func() (*string, error) {
			var status *string
			status, err := UserRepo.FindUserStatus(userID)
			return status, err
		})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domainerrors.ErrUnauthorized
	}
	return status, err
}
