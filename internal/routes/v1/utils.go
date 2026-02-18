package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/akashefrath/capecom-pm/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Utils(r *gin.RouterGroup, container container.Container) {
	utils := r.Group("utils")
	utils.Use(container.Middleware.Auth.VerifyToken(middleware.AllowedTypeAll))
	roleHandler := container.Handler.RoleHandler
	utils.GET("roles", roleHandler.GetAllActive)
}
