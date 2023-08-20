package handlers

import (
	"time"

	"auth/app/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

const (
	CSRF string = "csrf"
)

func CsrfMiddleware(c *fiber.Ctx) error {
	return csrf.New(csrf.Config{
		Storage:        store.Storage,
		KeyLookup:      "form:csrf",
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUID,
		ContextKey:     "csrf",
	})(c)
}
