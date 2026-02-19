package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/akashefrath/capecom-pm/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TimeClock(r *gin.RouterGroup, container container.Container) {
	timeClock := r.Group("time-clock")
	timeClock.Use(container.Middleware.Auth.VerifyToken(middleware.AllowedTypeAll))
	handler := container.Handler.TimeClockHandler
	timeClock.POST("clock-in", handler.ClockIn)
	timeClock.POST("clock-out", handler.ClockOut)

	timeClock.POST("break-in", handler.BreakIn)
	timeClock.POST("break-out", handler.BreakOut)

	timeClock.POST("time-out")
}
