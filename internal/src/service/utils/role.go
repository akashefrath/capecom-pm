package utilsservice

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	utilsrepository "github.com/akashefrath/capecom-pm/internal/src/repository/utils"
)

type Role struct {
	RoleRepo *utilsrepository.Role
}

func NewRole(roleRepo *utilsrepository.Role) *Role {
	return &Role{RoleRepo: roleRepo}
}

func (s *Role) GetAllActive() ([]dto.RoleResponse, error) {
	return s.RoleRepo.GetAllActive()
}
