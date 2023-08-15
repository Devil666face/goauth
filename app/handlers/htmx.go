package handlers

import (
	"github.com/gofiber/fiber/v2"
)

var (
	Htmx      string = "HTMX"
	HxRequest string = "Hx-Request"
)

func SetHtmx(app fiber.Router) {
	app.Use(HtmxMiddleware)
}

func HtmxMiddleware(c *fiber.Ctx) error {
	c.Locals(Htmx, false)
	if _, ok := c.GetReqHeaders()[HxRequest]; ok {
		c.Locals(Htmx, true)
	}
	return c.Next()
}
