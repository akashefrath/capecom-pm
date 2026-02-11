package container

import (
	"capecom-pm/internal/middleware"
	jwtutil "capecom-pm/internal/utils/jwt"
)

type Middleware struct {
	AdminMiddleware *middleware.AdminMiddleware
	UserMiddleware  *middleware.UserMiddleware
}

func NewMiddleware(jwtManager *jwtutil.Manager, repository *Repository) *Middleware {
	return &Middleware{
		AdminMiddleware: middleware.NewAdminMiddleware(jwtManager, repository.UserRepo, repository.CacheRepo,repository.SessionRepo),
		UserMiddleware:  middleware.NewUserMiddleware(jwtManager, repository.UserRepo, repository.CacheRepo,repository.SessionRepo),
	}
}
