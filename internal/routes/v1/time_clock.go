package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/akashefrath/capecom-pm/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TimeClock(r *gin.RouterGroup, container container.Container) {
	timeClock := r.Group("time-clock")
	timeClock.Use(container.Middleware.Auth.VerifyToken(middleware.AllowedTypeAll))
	handler := container.Handler.TimeClockHandler

	timeClock.GET("details", handler.GetTodayDetails)

	timeClock.POST("clock-in", func(c *gin.Context) {
		handler.AdvancePunch(c, models.LogIn)
	})
	timeClock.POST("clock-out", func(c *gin.Context) {
		handler.AdvancePunch(c, models.LogOut)
	})

	timeClock.POST("break-in", func(c *gin.Context) {
		handler.AdvancePunch(c, models.BrakeIn)
	})
	timeClock.POST("break-out", func(c *gin.Context) {
		handler.AdvancePunch(c, models.BrakeOut)
	})

}

///	timeClock.POST("time-out")
