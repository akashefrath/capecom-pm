package adminv1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/gin-gonic/gin"
)

func Users(r *gin.RouterGroup, container container.Container) {

	userGroup := r.Group("/users")
	userHandler := container.Handler.UserHandler
	userGroup.POST("", userHandler.Create)
	userGroup.GET("", userHandler.GetUsers)
	userGroup.GET(":id", userHandler.GetUserByID)
}
