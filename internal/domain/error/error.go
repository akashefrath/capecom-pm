package domainerrors

import "errors"

var (

	// ErrInvalidCredentials auth
	ErrInvalidCredentials = errors.New("invalid_login_credentials")
	ErrUnauthorized       = errors.New("unauthorized")

	// ErrUserNotFound user
	ErrUserNotFound   = errors.New("user_not_found")
	ErrDuplicateEmail = errors.New("duplicate_email")

	//ErrBadRequest common
	ErrBadRequest = errors.New("bad_request")
	ErrInternal   = errors.New("internal_error")
)
