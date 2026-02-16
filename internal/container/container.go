package container

import (
	"database/sql"

	"github.com/akashefrath/capecom-pm/internal/config"
)

type Container struct {
	DB           *sql.DB
	Config       *config.Config
	Handler      *Handler
	Service      *Service
	Repositories *Repositories
}

func New(db *sql.DB, config *config.Config) Container {
	repo := NewRepository(db)
	service := NewService(repo)

	return Container{
		DB:           db,
		Config:       config,
		Handler:      SetupHandler(service),
		Service:      service,
		Repositories: repo,
	}
}
