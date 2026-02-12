package i18n

var TA = map[string]string{

	// validation
	"required":     "இந்த புலம் அவசியம்",
	"email":        "தவறான மின்னஞ்சல் வடிவம்",
	"min":          "குறைந்தது %s எழுத்துகள் வேண்டும்",
	"max":          "அதிகபட்சம் %s எழுத்துகள் மட்டும்",
	"invalid_body": "தவறான கோரிக்கை உடல்",

	// common
	"bad_request":    "தவறான கோரிக்கை",
	"internal_error": "உள்ளக சேவையக பிழை",
	"unauthorized":   "அனுமதி இல்லை",

	// user
	"user_not_found":  "பயனர் கிடைக்கவில்லை",
	"duplicate_email": "மின்னஞ்சல் ஏற்கனவே உள்ளது",

	// auth
	"invalid_login_credentials": "தவறான உள்நுழைவு விவரங்கள்",

	// client
	"client_not_found": "வாடிக்கையாளர் கிடைக்கவில்லை",
	"duplicate_client": "வாடிக்கையாளர் ஏற்கனவே உள்ளது",
}
