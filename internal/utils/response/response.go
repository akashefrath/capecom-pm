package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, code int, res APIResponse) {
	c.JSON(code, res)

}
func JSONOk(c *gin.Context, res APIResponse) {
	JSON(c, http.StatusOK, res)
}
func JSONCreated(c *gin.Context, res APIResponse) {
	JSON(c, http.StatusCreated, res)
}
