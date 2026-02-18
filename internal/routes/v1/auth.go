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

	authSecure := auth.Group("")
	authSecure.Use(authMiddleware.VerifyToken(middleware.AllowedTypeAll))
	authSecure.GET("me", authHandler.Me)
	authSecure.POST("logout", authHandler.Logout)
}
