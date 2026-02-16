package response

import (
	"errors"
	"net/http"

	domainerrors "github.com/akashefrath/capecom-pm/internal/domain/error"
	"github.com/akashefrath/capecom-pm/internal/utils/i18n"
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
	case errors.Is(err, domainerrors.ErrClientNotFound):
		code = http.StatusNotFound
	case errors.Is(err, domainerrors.ErrDuplicateClient):
		code = http.StatusConflict
	case errors.Is(err, domainerrors.ErrProjectNotFound):
		code = http.StatusNotFound
	case errors.Is(err, domainerrors.ErrDuplicateProject):
		code = http.StatusConflict
	case errors.Is(err, domainerrors.ErrProjectAssetNotFound):
		code = http.StatusNotFound
	case errors.Is(err, domainerrors.ErrTicketNotFound):
		code = http.StatusNotFound
	case errors.Is(err, domainerrors.ErrDuplicateTicket):
		code = http.StatusConflict
	case errors.Is(err, domainerrors.ErrTicketTypeNotFound):
		code = http.StatusBadRequest
	case errors.Is(err, domainerrors.ErrNotProjectMember):
		code = http.StatusBadRequest
	case errors.Is(err, domainerrors.ErrNotTicketAssignee):
		code = http.StatusForbidden
	case errors.Is(err, domainerrors.ErrTimeEntryNotFound):
		code = http.StatusNotFound
	case errors.Is(err, domainerrors.ErrFileNotFound):
		code = http.StatusBadRequest
	case errors.Is(err, domainerrors.ErrFileNotUploaded):
		code = http.StatusBadRequest
	case errors.Is(err, domainerrors.ErrBadRequest):
		code = http.StatusBadRequest
	case errors.Is(err, domainerrors.ErrUnauthorized):
		code = http.StatusUnauthorized
	}

	message := i18n.GetMessages(lang)
	finalMessage := message[err.Error()]
	if finalMessage == "" {
		finalMessage = err.Error()
	}
	var appErr *domainerrors.AppError
	var entity string
	var funcErr string
	if errors.As(err, &appErr) {
		if appErr.Code() != nil {
			code = *appErr.Code()
			entity = appErr.Entity()
			funcErr = appErr.Function()

			//finalMessage = message[appErr.Error()]

		}
	}
	JSON(c, code, APIResponse{
		Success: false,
		Message: finalMessage,
		Entity:  entity,
		Func:    funcErr,
	})
}
