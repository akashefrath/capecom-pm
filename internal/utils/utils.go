package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) string {
	return c.GetString("userID")

}
func ToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
func ToString(v any) string {
	return fmt.Sprint(v)
}
