package container

import (
	"database/sql"

	"github.com/akashefrath/capecom-pm/internal/src/repository"
)

type Repositories struct {
	Auth *repository.Auth
	User *repository.User
}

func NewRepository(db *sql.DB) *Repositories {
	return &Repositories{
		Auth: repository.NewAuth(db),
		User: repository.NewUser(db),
	}
}
