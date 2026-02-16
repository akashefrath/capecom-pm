package i18n

import "strings"

func GetMessages(lang string) map[string]string {
	lang = strings.ToLower(lang)

	switch lang {
	case "ta":
		return TA
	default:
		return EN
	}
}
