package adminv1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/akashefrath/capecom-pm/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Admin(r *gin.RouterGroup, container container.Container) {
	adminGroup := r.Group("/admin")
	adminGroup.Use(container.Middleware.Auth.VerifyToken(middleware.AllowedTypeAdmin))
	Users(adminGroup, container)
}
