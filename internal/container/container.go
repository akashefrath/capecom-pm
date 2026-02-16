package container

import (
	"github.com/akashefrath/capecom-pm/internal/config"
	jwtutil "github.com/akashefrath/capecom-pm/internal/utils/jwt"
	"github.com/jmoiron/sqlx"
)

type Container struct {
	DB           *sqlx.DB
	Config       *config.Config
	Handler      *Handler
	Service      *Service
	Repositories *Repositories
	JWTManager   *jwtutil.Manager
}

func New(db *sqlx.DB, config *config.Config) Container {
	jwtManager := jwtutil.NewJWTManager(
		config.JWT.UserSecret,
		config.JWT.UserRefreshSecret,
		config.JWT.AdminSecret,
		config.JWT.AdminRefreshSecret,
		config.JWT.ExpireHours,
		config.JWT.RefreshExpireHours,
	)
	repo := NewRepository(db, config)
	service := NewService(repo, jwtManager)

	return Container{
		DB:           db,
		Config:       config,
		Handler:      SetupHandler(service),
		Service:      service,
		Repositories: repo,
		JWTManager:   jwtManager,
	}
}
