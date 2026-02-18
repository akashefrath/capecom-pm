package utilshandler

import (
	utilsservice "github.com/akashefrath/capecom-pm/internal/src/service/utils"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type UtilsHandler struct {
	Utils *utilsservice.Utils
}

func NewUtilsHandler(utils *utilsservice.Utils) UtilsHandler {
	return UtilsHandler{Utils: utils}
}

func (h UtilsHandler) GetRoles(c *gin.Context) {
	response.JSONOk(c, response.APIResponse{})

}
