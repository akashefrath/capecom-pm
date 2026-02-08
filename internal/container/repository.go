package container

import (
	"capecom-pm/internal/repositories"
	mastersreo "capecom-pm/internal/repositories/masters"

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
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AuthRepo:        repositories.NewAuthRepo(db),
		UserRepo:        repositories.NewUserRepo(db),
		RoleRepo:        mastersreo.NewRoleRepo(db),
		GroupRepo:       mastersreo.NewGroupRepo(db),
		DesignationRepo: mastersreo.NewDesignationRepo(db),
		DepartmentRepo:  mastersreo.NewDepartmentRepo(db),
		MasterDataRepo:  mastersreo.NewMasterDataRepo(db),
	}
}
