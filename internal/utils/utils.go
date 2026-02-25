package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/akashefrath/capecom-pm/internal/utils/i18n"
	"github.com/gin-gonic/gin"
)

func GetUserUuid(c *gin.Context) string {
	return c.GetString(CtxUserUUID)
}
func GetUserID(c *gin.Context) int64 {
	return c.GetInt64(CtxUserID)
}
func GetIsAdmin(c *gin.Context) bool {
	return c.GetBool(CtxIsAdmin)
}

func GetJTI(c *gin.Context) string {
	return c.GetString(CtxJTI)
}

func GetUserData(c *gin.Context) (string, int64, bool) {
	return GetUserUuid(c), GetUserID(c), GetIsAdmin(c)
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
func GetMessageWithExtra(key string, c *gin.Context, a ...string) string {
	lang := c.GetHeader("Accept-Language")
	message := i18n.GetMessages(lang)
	finalMessage := message[key]
	for _, v := range a {
		trMessage := message[v]
		if trMessage == "" {
			trMessage = v
		}
		finalMessage = fmt.Sprintf(finalMessage, trMessage)
	}
	return finalMessage
}

func ParseDate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil
	}
	return &t
}
func CombineDateTime(baseDate time.Time, userTimeStr string) (time.Time, error) {
	// Define supported layouts (24h and 12h)
	layouts := []string{"15:04", "3:04 PM", "3:04PM", "15.04"}

	var parsedTime time.Time
	var err error

	// Try parsing the time string against the layouts
	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, userTimeStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %s", userTimeStr)
	}

	// Combine the fixed date with the parsed time
	combined := time.Date(
		baseDate.Year(), baseDate.Month(), baseDate.Day(),
		parsedTime.Hour(), parsedTime.Minute(), 0, 0, baseDate.Location(),
	)

	return combined, nil
}
