package middlewares

import (
	"fmt"

	"app/config"

	"github.com/gofiber/fiber/v2"
)

const (
	PROTO string = "https"
)

func HttpsRedirectMiddleware(c *fiber.Ctx) error {
	return c.Redirect(fmt.Sprintf("%s://%s:%s", PROTO, config.ALLOW_HOST, config.HTTPS_PORT))
}
