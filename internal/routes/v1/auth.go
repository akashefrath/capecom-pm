package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/akashefrath/capecom-pm/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Auth(r *gin.RouterGroup, container container.Container) {
	auth := r.Group("auth")
	authHandler := container.Handler.AuthHandler
	auth.POST("login", authHandler.Login)
	auth.POST("refresh-token", authHandler.RefreshToken)
	authMiddleware := container.Middleware.Auth
	auth.Use(authMiddleware.VerifyToken(middleware.AllowedTypeAll))
	auth.GET("me", authHandler.Me)
	auth.POST("logout", authHandler.Logout)
}
