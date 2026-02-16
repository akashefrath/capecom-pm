package container

import (
	"github.com/akashefrath/capecom-pm/internal/config"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	Session *repository.Session
	User    *repository.User
}

func NewRepository(db *sqlx.DB, config *config.Config) *Repositories {
	return &Repositories{
		Session: repository.NewSession(db, config),
		User:    repository.NewUser(db),
	}
}
