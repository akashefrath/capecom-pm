package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/akashefrath/capecom-pm/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AttendanceSummary(r *gin.RouterGroup, container container.Container) {
	summary := r.Group("summary")
	summary.Use(container.Middleware.Auth.VerifyToken(middleware.AllowedTypeAll))
	handler := container.Handler.AttendanceSummaryHandler

	summary.GET("list", handler.GetList)
	summary.GET(":uuid", handler.GetByUUID)
}
