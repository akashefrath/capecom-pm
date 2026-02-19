package container

import (
	"github.com/akashefrath/capecom-pm/internal/src/service"
	utilsservice "github.com/akashefrath/capecom-pm/internal/src/service/utils"
	jwtutil "github.com/akashefrath/capecom-pm/internal/utils/jwt"
)

type Service struct {
	Auth                  *service.Auth
	User                  *service.User
	Role                  *utilsservice.Role
	Utils                 *utilsservice.Utils
	AttendancePolicy      *service.AttendancePolicy
	AttendancePolicyGroup *service.AttendancePolicyGroup
	ShiftSystem           *service.ShiftSystem
	ShiftSystemGroup      *service.ShiftSystemGroup
}

func NewService(repo *Repositories, jwtManager *jwtutil.Manager) *Service {
	return &Service{
		Auth:                  service.NewAuth(repo.Session, repo.User, jwtManager, repo.Redis),
		User:                  service.NewUser(repo.User),
		Role:                  utilsservice.NewRole(repo.Role),
		Utils:                 utilsservice.NewUtils(repo.Utils),
		AttendancePolicy:      service.NewAttendancePolicy(repo.AttendancePolicy),
		AttendancePolicyGroup: service.NewAttendancePolicyGroup(repo.AttendancePolicyGroup),
		ShiftSystem:           service.NewShiftSystem(repo.ShiftSystem),
		ShiftSystemGroup:      service.NewShiftSystemGroup(repo.ShiftSystemGroup, repo.ShiftSystem),
	}
}
