package middlewares

import (
	"strings"

	"app/config"

	"github.com/gofiber/fiber/v2"
)

const (
	HOST string = "Host"
)

func AllowedHostMiddleware(c *fiber.Ctx) error {
	if host, ok := c.GetReqHeaders()[HOST]; ok {
		if strings.Contains(host, config.ALLOW_HOST) {
			return c.Next()
		}
	}
	return fiber.ErrBadRequest
}
