package response

import "github.com/gin-gonic/gin"

func JSON(c *gin.Context, code int, res APIResponse) {
	c.JSON(code, res)
}
