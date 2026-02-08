package response

import (
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/utils/i18n"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FromError(c *gin.Context, err error) {
	lang := c.GetHeader("Accept-Language")
	code := http.StatusInternalServerError
	switch {
	case errors.Is(err, domainerrors.ErrInvalidCredentials):
		code = http.StatusUnauthorized
	case errors.Is(err, domainerrors.ErrDuplicateEmail):
		code = http.StatusConflict
	case errors.Is(err, domainerrors.ErrUserNotFound):
		code = http.StatusNotFound
	case errors.Is(err, domainerrors.ErrBadRequest):
		code = http.StatusBadRequest
	case errors.Is(err, domainerrors.ErrUnauthorized):
		code = http.StatusUnauthorized
	}

	message := i18n.GetMessages(lang)
	println(err.Error())
	JSON(c, code, APIResponse{
		Success: false,
		Message: message[err.Error()],
	})
}
