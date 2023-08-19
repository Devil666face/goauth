package routes

import (
	"auth/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func FreeRoutes(app fiber.Router) {
	r := app.Group("/auth")
	r.Get("/login", handlers.LoginGet).Name("login")
	r.Post("/login", handlers.LoginPost)
	r.Get("/logout", handlers.LogoutGet).Name("logout")
	r.Get("/health", handlers.Health)
}

func SuperUserRoutes(app fiber.Router) {
	r := app.Group("/auth")
	r.Use(handlers.AuthMiddleware)
	r.Use(handlers.SuperUserMiddleware)
	r.Use(handlers.HtmxMiddleware)

	r.Get("/users", handlers.UserControlGet).Name("users")

	u := r.Group("/user")
	u.Get("/new", handlers.UserCreateGet).Name("user-new")
	u.Post("/new", handlers.UserCreatePost)
	u.Get("/:id<int>/edit", handlers.UserEditGet).Name("user-edit")
	u.Post("/:id<int>/edit", handlers.UserEditPost)
	u.Delete("/:id<int>/delete", handlers.UserDeletePost).Name("user-delete")
}

func AuthRoutes(app fiber.Router) {
	r := app.Group("")
	r.Get("/user", handlers.AuthMiddleware, handlers.UserGet)
}
