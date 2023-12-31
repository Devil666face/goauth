package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

const (
	HTMX      string = "HTMX"
	HXREQUEST string = "Hx-Request"
)

func HtmxMiddleware(c *fiber.Ctx) error {
	c.Locals(HTMX, false)
	if _, ok := c.GetReqHeaders()[HXREQUEST]; ok {
		c.Locals(HTMX, true)
	}
	return c.Next()
}

// [HTMX Middleware]
// if c.Locals(HTMX).(bool) {
// 	return c.Status(fiber.StatusOK).SendString("<div>Missmatch username or password</biv>")
// }
