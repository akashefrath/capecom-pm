package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/gin-gonic/gin"
)

func Auth(r *gin.RouterGroup, container container.Container) {
	auth := r.Group("auth")
	authHandler := container.Handler.AuthHandler
	auth.POST("login", authHandler.Login)

}
