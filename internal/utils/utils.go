package utils

import (
	"capecom-pm/internal/utils/i18n"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) string {
	return c.GetString("userID")
}
func GetJTI(c *gin.Context) string {
	return c.GetString("jti")
}

func ToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
func ToString(v any) string {
	return fmt.Sprint(v)
}

func GetMessage(key string, c *gin.Context) string {
	lang := c.GetHeader("Accept-Language")
	message := i18n.GetMessages(lang)
	finalMessage := message[key]
	return finalMessage
}
