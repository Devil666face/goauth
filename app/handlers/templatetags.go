package handlers

import (
	"auth/app/utils"

	"github.com/gofiber/fiber/v2"
)

func TemplatetagMiddleware(c *fiber.Ctx) error {
	utils.Ctx = c
	return c.Next()
}
