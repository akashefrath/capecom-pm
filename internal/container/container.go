package container

import (
	"github.com/akashefrath/capecom-pm/internal/config"
	"github.com/akashefrath/capecom-pm/internal/database"
	"github.com/akashefrath/capecom-pm/internal/middleware"
	jwtutil "github.com/akashefrath/capecom-pm/internal/utils/jwt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	DB           *sqlx.DB
	DBtx         *database.Database
	Config       *config.Config
	Handler      *Handler
	Service      *Service
	Repositories *Repositories
	JWTManager   *jwtutil.Manager
	Middleware   *Middleware
}

type Middleware struct {
	Auth *middleware.Auth
}

func New(db *sqlx.DB, config *config.Config, redis *redis.Client) Container {
	jwtManager := jwtutil.NewJWTManager(
		config.JWT.UserSecret,
		config.JWT.UserRefreshSecret,
		config.JWT.AdminSecret,
		config.JWT.AdminRefreshSecret,
		config.JWT.ExpireHours,
		config.JWT.RefreshExpireHours,
	)
	dbTX := &database.Database{DB: db}
	repo := NewRepository(db, config, redis, dbTX)
	service := NewService(repo, jwtManager)

	return Container{
		DB:           db,
		DBtx:         dbTX,
		Config:       config,
		Handler:      SetupHandler(service),
		Service:      service,
		Repositories: repo,
		JWTManager:   jwtManager,
		Middleware:   SetMiddleware(jwtManager, repo),
	}
}

func SetMiddleware(jwtManager *jwtutil.Manager, repo *Repositories) *Middleware {
	return &Middleware{Auth: middleware.NewAuth(jwtManager, repo.User, repo.Redis, repo.Session)}
}
