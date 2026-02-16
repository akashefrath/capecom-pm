package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, container container.Container) {
	v1API := r.Group("api/v1")
	Auth(v1API, container)

}
