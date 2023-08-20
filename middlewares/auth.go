package middlewares

import (
	"fmt"

	"app/models"
	"app/store"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	session, err := store.Store.Get(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("login", fiber.Map{})
	}

	if session.Get(store.AUTH_KEY) == nil {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("login", fiber.Map{})
	}

	uid := session.Get(store.USER_ID)
	if uid == nil {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("login", fiber.Map{})
	}

	u := new(models.User)
	if models.GetUser(u, fmt.Sprint(uid)); err != nil {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("login", fiber.Map{})
	}

	c.Locals(models.USER, u)

	return c.Next()
}

func SuperUserMiddleware(c *fiber.Ctx) error {
	u := c.Locals(models.USER)
	user, ok := u.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).RedirectToRoute("login", fiber.Map{})
	}
	if !user.Admin {
		return fiber.ErrNotFound
	}
	return c.Next()
}
