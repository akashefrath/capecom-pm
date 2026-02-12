package container

import (
	"capecom-pm/internal/services"
	"capecom-pm/internal/storage"
	jwtutil "capecom-pm/internal/utils/jwt"
)

type Service struct {
	AuthService    *services.AuthService
	UserService    *services.UserService
	FileService    *services.FileService
	ClientService  *services.ClientService
	ProjectService *services.ProjectService
	UtilsService   *services.UtilsService
}

func NewService(repository *Repository, jwt *jwtutil.Manager, r2Client *storage.R2Client) *Service {
	return &Service{
		AuthService: services.NewAuthService(repository.AuthRepo, jwt, repository.UserRepo, repository.SessionRepo, repository.CacheRepo),
		UserService: services.NewUserService(
			repository.UserRepo,
			repository.MasterDataRepo,
		),
		FileService:    services.NewFileService(repository.FileRepo, repository.UserRepo, r2Client, repository.CacheRepo),
		ClientService:  services.NewClientService(repository.ClientRepo, repository.UserRepo, repository.CacheRepo),
		ProjectService: services.NewProjectService(repository.ProjectRepo, repository.ClientRepo, repository.UserRepo, repository.CacheRepo),
		UtilsService:   services.NewUtilsService(repository.UtilsRepo),
	}
}
