package container

import (
	"github.com/akashefrath/capecom-pm/internal/config"
	"github.com/akashefrath/capecom-pm/internal/database"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
	utilsrepository "github.com/akashefrath/capecom-pm/internal/src/repository/utils"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	Session *repository.Session
	User    *repository.User
	Role    *utilsrepository.Role
	Redis   *utilsrepository.Redis
	Utils   *utilsrepository.Utils
}

func NewRepository(db *sqlx.DB, config *config.Config, redis *redis.Client, dbTX *database.Database) *Repositories {
	return &Repositories{
		Session: repository.NewSession(db, config),
		User:    repository.NewUser(db, dbTX),
		Role:    utilsrepository.NewRole(db),
		Redis:   utilsrepository.NewRedis(redis),
		Utils:   utilsrepository.NewUtils(db),
	}
}
