package container

import (
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	mastersreo "capecom-pm/internal/repositories/masters"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	AuthRepo        *repositories.AuthRepo
	UserRepo        *repositories.UserRepo
	RoleRepo        *mastersreo.RoleRepo
	GroupRepo       *mastersreo.GroupRepo
	DesignationRepo *mastersreo.DesignationRepo
	DepartmentRepo  *mastersreo.DepartmentRepo
	MasterDataRepo  *mastersreo.MasterDataRepo
	CacheRepo       *cacherepo.RedisRepo
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		AuthRepo:        repositories.NewAuthRepo(db),
		UserRepo:        repositories.NewUserRepo(db),
		RoleRepo:        mastersreo.NewRoleRepo(db),
		GroupRepo:       mastersreo.NewGroupRepo(db),
		DesignationRepo: mastersreo.NewDesignationRepo(db),
		DepartmentRepo:  mastersreo.NewDepartmentRepo(db),
		MasterDataRepo:  mastersreo.NewMasterDataRepo(db),
		CacheRepo:       cacherepo.NewRedisRepo(redis),
	}
}
