package container

import (
	"capecom-pm/internal/middleware"
	jwtutil "capecom-pm/internal/utils/jwt"
)

type Middleware struct {
	AdminMiddleware *middleware.AdminMiddleware
}

func NewMiddleware(jwtManager *jwtutil.JWTManager, repository *Repository) *Middleware {
	return &Middleware{
		AdminMiddleware: middleware.NewAdminMiddleware(jwtManager, repository.UserRepo),
	}
}
