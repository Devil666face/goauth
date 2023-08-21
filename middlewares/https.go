package middlewares

import (
	"fmt"
	// "strings"

	"app/config"

	"github.com/gofiber/fiber/v2"
)

const (
	HTTP string = "http"
)

func HttpsRedirectMiddleware(c *fiber.Ctx) error {
	if config.TLS != "True" {
		return c.Next()
	}
	if c.Protocol() == HTTP {
		return c.Redirect(fmt.Sprintf("https://%s", config.CONNECT_HTTPS))
	}
	return c.Next()
}
