package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SetHtmx(app fiber.Router) {
	app.Use(HtmxMiddleware)
}

func HtmxMiddleware(c *fiber.Ctx) error {
	c.Locals("HTMX", false)
	if _, ok := c.GetReqHeaders()["Hx-Request"]; ok {
		c.Locals("HTMX", true)
	}
	fmt.Println(c.Locals("HTMX"))
	return c.Next()
}
