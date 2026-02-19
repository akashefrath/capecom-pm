package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/akashefrath/capecom-pm/internal/routes/v1/admin"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, container container.Container) {
	v1API := r.Group("api/v1")

	adminv1.Admin(v1API, container)
	Auth(v1API, container)
	Utils(v1API, container)
	TimeClock(v1API, container)

}
