package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Check if it's an fiber.Error type
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Render("error", fiber.Map{
		"Statuscode": code,
		"Error":      err.Error(),
		"Title":      fmt.Sprintf("Error %d", code),
	})
}
