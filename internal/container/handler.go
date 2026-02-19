package container

import (
	"github.com/akashefrath/capecom-pm/internal/src/handler"
	utilshandler "github.com/akashefrath/capecom-pm/internal/src/handler/utils"
)

type Handler struct {
	AuthHandler                  handler.AuthHandler
	UserHandler                  handler.UserHandler
	RoleHandler                  handler.RoleHandler
	UtilsHandler                 utilshandler.UtilsHandler
	AttendancePolicyHandler      handler.AttendancePolicyHandler
	AttendancePolicyGroupHandler handler.AttendancePolicyGroupHandler
	ShiftSystemHandler           handler.ShiftSystemHandler
	ShiftSystemGroupHandler      handler.ShiftSystemGroupHandler
}

func SetupHandler(service *Service) *Handler {
	return &Handler{
		AuthHandler:                  handler.NewAuth(service.Auth),
		UserHandler:                  handler.NewUser(service.User),
		RoleHandler:                  handler.NewRole(service.Role),
		UtilsHandler:                 utilshandler.NewUtilsHandler(service.Utils),
		AttendancePolicyHandler:      handler.NewAttendancePolicy(service.AttendancePolicy),
		AttendancePolicyGroupHandler: handler.NewAttendancePolicyGroup(service.AttendancePolicyGroup),
		ShiftSystemHandler:           handler.NewShiftSystem(service.ShiftSystem),
		ShiftSystemGroupHandler:      handler.NewShiftSystemGroup(service.ShiftSystemGroup),
	}
}
