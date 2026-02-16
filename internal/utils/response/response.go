package response

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func JSON(c fiber.Ctx, code int, res APIResponse) error {
	err := c.Status(code).JSON(res)
	if err != nil {
		return err
	}

	return nil
}
func JSONOk(c fiber.Ctx, res APIResponse) error {
	return JSON(c, http.StatusOK, res)
}
func JSONCreated(c fiber.Ctx, res APIResponse) error {
	return JSON(c, http.StatusCreated, res)
}
