package bind

import (
	"capecom-pm/internal/utils/i18n"
	"capecom-pm/internal/utils/response"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldErrors map[string][]string

func AndValidate(c *gin.Context, req any, entity string) bool {

	if err := c.ShouldBind(req); err != nil {

		var ve validator.ValidationErrors

		lang := c.GetHeader("Accept-Language")
		msgs := i18n.GetMessages(lang)

		errs := FieldErrors{}

		if errors.As(err, &ve) {
			for _, fe := range ve {

				field := strings.ToLower(fe.Field())
				tag := fe.Tag()
				param := fe.Param()

				msg := msgs[tag]
				if msg == "" {
					msg = tag
				}

				if param != "" {
					msg = fmt.Sprintf(msg, param)
				}

				errs[field] = append(errs[field], msg)
			}
		} else {
			errs["body"] = []string{msgs["invalid_body"]}
		}

		response.JSON(c, http.StatusBadRequest, response.APIResponse{
			Success: false,
			Errors:  errs,
			Func:    "validate",
			Entity:  entity,
		})

		return false
	}

	return true
}

func QueryBinder(c *gin.Context, req any, entity string) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return err
	}

	return nil
}
