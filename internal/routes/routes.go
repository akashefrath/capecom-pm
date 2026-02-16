package routes

import (
	"net/http"

	"github.com/akashefrath/capecom-pm/internal/container"
	v1 "github.com/akashefrath/capecom-pm/internal/routes/v1"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, container container.Container) {
	//	app.Use(r.New())
	r.Any("/ping", func(c *gin.Context) {
		response.JSON(c, http.StatusOK, response.APIResponse{
			Success: true,
			Message: "Pong",
		})
	})

	v1.Setup(r, container)
}
