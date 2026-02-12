package container

import (
	"capecom-pm/internal/config"
	"capecom-pm/internal/storage"
	jwtutil "capecom-pm/internal/utils/jwt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	Handler     *Handler
	Service     *Service
	Repository  *Repository
	Middleware  *Middleware
	JWTManager  *jwtutil.Manager
	RedisClient *redis.Client
	R2Client    *storage.R2Client
}

func NewContainer(db *gorm.DB, cfg config.Config, redis *redis.Client) *Container {

	jwtManager := jwtutil.NewJWTManager(
		cfg.JWT.UserSecret,
		cfg.JWT.UserRefreshSecret,
		cfg.JWT.AdminSecret,
		cfg.JWT.AdminRefreshSecret,
		cfg.JWT.ExpireHours,
		cfg.JWT.RefreshExpireHours,
	)

	r2Client := storage.NewR2Client(
		cfg.R2.AccountID,
		cfg.R2.AccessKeyID,
		cfg.R2.AccessKeySecret,
		cfg.R2.BucketName,
		cfg.R2.FolderName,
	)

	repository := NewRepository(db, redis)

	service := NewService(repository, jwtManager, r2Client)
	handler := NewHandler(service)
	middleware := NewMiddleware(jwtManager, repository)
	return &Container{
		Handler:     handler,
		Service:     service,
		Repository:  repository,
		Middleware:  middleware,
		JWTManager:  jwtManager,
		RedisClient: redis,
		R2Client:    r2Client,
	}
}
