package handlers

import (
	"strings"

	"auth/app/config"

	"github.com/gofiber/fiber/v2"
)

const (
	HOST string = "Host"
)

func AllowedHostMiddleware(c *fiber.Ctx) error {
	host := c.GetReqHeaders()[HOST]
	if strings.Contains(host, config.ALLOWED) {
		return c.Next()
	}
	return fiber.ErrBadRequest
}
