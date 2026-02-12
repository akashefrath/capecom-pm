package i18n

var EN = map[string]string{

	// validation
	"required":     "This field is required",
	"email":        "Invalid email format",
	"min":          "Minimum %s characters required",
	"max":          "Maximum %s characters allowed",
	"invalid_body": "Invalid request body",

	// common
	"bad_request":    "Bad request",
	"internal_error": "Internal server error",
	"unauthorized":   "Unauthorized access",

	// user
	"user_not_found":  "User not found",
	"duplicate_email": "Email already exists",

	// auth
	"invalid_login_credentials": "Invalid login credentials",

	// client
	"client_not_found": "Client not found",
	"duplicate_client": "Client already exists",

	"logged_out_success": "Logged out successfully",
}
