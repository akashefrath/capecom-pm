package container

import (
	"github.com/akashefrath/capecom-pm/internal/config"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
	utilsrepository "github.com/akashefrath/capecom-pm/internal/src/repository/utils"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	Session *repository.Session
	User    *repository.User
	Redis   *utilsrepository.Redis
}

func NewRepository(db *sqlx.DB, config *config.Config, redis *redis.Client) *Repositories {
	return &Repositories{
		Session: repository.NewSession(db, config),
		User:    repository.NewUser(db),
		Redis:   utilsrepository.NewRedis(redis),
	}
}
