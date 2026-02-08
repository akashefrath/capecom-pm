package utils

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) string {
	return c.GetString("userID")

}
