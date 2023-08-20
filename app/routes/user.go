package routes

import (
	"auth/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func FreeRoutes(app fiber.Router) {
	auth := app.Group("/auth")
	auth.Use(
		handlers.HtmxMiddleware,
		handlers.TemplatetagMiddleware,
	)
	auth.Get("/login", handlers.LoginGet).Name("login")
	auth.Post("/login", handlers.LoginPost)
	auth.Get("/logout", handlers.LogoutGet).Name("logout")
	auth.Get("/health", handlers.Health)
}

func SuperUserRoutes(app fiber.Router) {
	auth := app.Group("/auth")
	auth.Use(
		handlers.AuthMiddleware,
		handlers.SuperUserMiddleware,
		handlers.HtmxMiddleware,
		handlers.TemplatetagMiddleware,
	)

	auth.Get("/users", handlers.UserControlGet).Name("users")

	user := auth.Group("/user")
	user.Get("/new", handlers.UserCreateGet).Name("user-new")
	user.Post("/new", handlers.UserCreatePost)
	user.Get("/:id<int>/edit", handlers.UserEditGet).Name("user-edit")
	user.Post("/:id<int>/edit", handlers.UserEditPost)
	user.Delete("/:id<int>/delete", handlers.UserDeletePost).Name("user-delete")
}

func AuthRoutes(app fiber.Router) {
	free := app.Group("")
	free.Get("/user", handlers.AuthMiddleware, handlers.UserGet)
}
