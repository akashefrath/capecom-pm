package container

import "gorm.io/gorm"

type Container struct {
	Handler    *Handler
	Service    *Service
	Repository *Repository
}

func NewContainer(db *gorm.DB) *Container {

	repository := NewRepository(db)
	service := NewService(db, repository)
	handler := NewHandler(service)
	return &Container{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}
