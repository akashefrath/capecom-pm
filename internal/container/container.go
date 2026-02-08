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
	JWTManager *jwtutil.JWTManager
}

func NewContainer(db *gorm.DB, cfg config.Config) *Container {

	jwtManager := jwtutil.NewJWTManager(
		cfg.JWT.UserSecret,
		cfg.JWT.AdminSecret,
		cfg.JWT.ExpireHours,
		cfg.JWT.RefreshExpireHours,
	)

	repository := NewRepository(db)
	service := NewService(db, repository, jwtManager)
	handler := NewHandler(service)
	return &Container{
		Handler:    handler,
		Service:    service,
		Repository: repository,
		JWTManager: jwtManager,
	}
}
