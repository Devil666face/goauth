package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func SetHtmx(app fiber.Router) {
	app.Use(HtmxMiddleware)
}

var (
	Htmx      string = "HTMX"
	HxRequest string = "Hx-Request"
)

func HtmxMiddleware(c *fiber.Ctx) error {
	c.Locals(Htmx, false)
	if _, ok := c.GetReqHeaders()[HxRequest]; ok {
		c.Locals(Htmx, true)
	}
	return c.Next()
}
