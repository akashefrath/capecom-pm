package container

import (
	"capecom-pm/internal/middleware"
	jwtutil "capecom-pm/internal/utils/jwt"
)

type Middleware struct {
	AdminMiddleware *middleware.AdminMiddleware
	UserMiddleware  *middleware.UserMiddleware
	AuthMiddleware  *middleware.AuthMiddleware
	RABCMiddleware  *middleware.RABCMiddleware
}

func NewMiddleware(jwtManager *jwtutil.Manager, repository *Repository) *Middleware {
	return &Middleware{
		AdminMiddleware: middleware.NewAdminMiddleware(jwtManager, repository.UserRepo, repository.CacheRepo, repository.SessionRepo),
		UserMiddleware:  middleware.NewUserMiddleware(jwtManager, repository.UserRepo, repository.CacheRepo, repository.SessionRepo),
		AuthMiddleware:  middleware.NewAuthMiddleware(jwtManager, repository.UserRepo, repository.CacheRepo, repository.SessionRepo),
		RABCMiddleware:  middleware.NewRABCMiddleware(jwtManager, repository.UserRepo, repository.CacheRepo, repository.SessionRepo),
	}
}
