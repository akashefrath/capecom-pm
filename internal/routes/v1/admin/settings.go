package adminv1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/gin-gonic/gin"
)

func Settings(r *gin.RouterGroup, container container.Container) {
	settings := r.Group("settings")
	attendancePolicyHandler := container.Handler.AttendancePolicyHandler
	attendancePolicyGroupHandler := container.Handler.AttendancePolicyGroupHandler
	shiftSystemHandler := container.Handler.ShiftSystemHandler
	shiftSystemGroupHandler := container.Handler.ShiftSystemGroupHandler

	attendancePolicies := settings.Group("attendance-policies")
	attendancePolicies.POST("", attendancePolicyHandler.Create)
	attendancePolicies.PUT("/:uuid", attendancePolicyHandler.Update)
	attendancePolicies.PUT("/:uuid/set-default", attendancePolicyHandler.SetDefault)
	attendancePolicies.DELETE("/:uuid", attendancePolicyHandler.Delete)
	attendancePolicies.GET("", attendancePolicyHandler.GetAll)
	attendancePolicies.GET("utils", attendancePolicyHandler.GetAllUtils)
	attendancePolicies.GET("/:uuid", attendancePolicyHandler.GetByUUID)

	attendancePolicyGroups := settings.Group("attendance-policy-groups")
	attendancePolicyGroups.POST("", attendancePolicyGroupHandler.Create)
	attendancePolicyGroups.PUT("/:uuid", attendancePolicyGroupHandler.Update)
	attendancePolicyGroups.DELETE("/:uuid", attendancePolicyGroupHandler.Delete)
	attendancePolicyGroups.GET("", attendancePolicyGroupHandler.GetAll)
	attendancePolicyGroups.GET("/:uuid", attendancePolicyGroupHandler.GetByUUID)

	shiftSystems := settings.Group("shift-systems")
	shiftSystems.POST("", shiftSystemHandler.Create)
	shiftSystems.PUT("/:uuid", shiftSystemHandler.Update)
	shiftSystems.PUT("/:uuid/set-default", shiftSystemHandler.SetDefault)
	shiftSystems.DELETE("/:uuid", shiftSystemHandler.Delete)
	shiftSystems.GET("", shiftSystemHandler.GetAll)
	shiftSystems.GET("utils", shiftSystemHandler.GetAllUtils)
	shiftSystems.GET("/:uuid", shiftSystemHandler.GetByUUID)

	shiftSystemGroups := settings.Group("shift-system-groups")
	shiftSystemGroups.POST("", shiftSystemGroupHandler.Create)
	shiftSystemGroups.PUT("/:uuid", shiftSystemGroupHandler.Update)
	shiftSystemGroups.DELETE("/:uuid", shiftSystemGroupHandler.Delete)
	shiftSystemGroups.GET("", shiftSystemGroupHandler.GetAll)
	shiftSystemGroups.GET("/:uuid", shiftSystemGroupHandler.GetByUUID)
}
