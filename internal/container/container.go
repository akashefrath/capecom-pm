package container

import (
	"capecom-pm/internal/config"
	jwtutil "capecom-pm/internal/utils/jwt"

	"gorm.io/gorm"
)

type Container struct {
	Handler    *Handler
	Service    *Service
	Repository *Repository
	Middleware *Middleware
	JWTManager *jwtutil.JWTManager
}

func NewContainer(db *gorm.DB, cfg config.Config) *Container {

	jwtManager := jwtutil.NewJWTManager(
		cfg.JWT.UserSecret,
		cfg.JWT.UserRefreshSecret,
		cfg.JWT.AdminSecret,
		cfg.JWT.AdminRefreshSecret,
		cfg.JWT.ExpireHours,
		cfg.JWT.RefreshExpireHours,
	)

	repository := NewRepository(db)
	service := NewService(db, repository, jwtManager)
	handler := NewHandler(service)
	middleware := NewMiddleware(jwtManager, repository)
	return &Container{
		Handler:    handler,
		Service:    service,
		Repository: repository,
		Middleware: middleware,
		JWTManager: jwtManager,
	}
}
