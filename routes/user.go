package routes

import (
	"app/handlers"

	. "app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func FreeRoutes(app fiber.Router) {
	auth := app.Group("/auth")
	auth.Get("/login", handlers.LoginGet).Name("login")
	auth.Post("/login", handlers.LoginPost)
	auth.Get("/logout", handlers.LogoutGet).Name("logout")
	auth.Get("/health", handlers.Health)
}

func SuperUserRoutes(app fiber.Router) {
	auth := app.Group("/auth")
	auth.Use(
		AuthMiddleware,
		SuperUserMiddleware,
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
	app.Use(AuthMiddleware)
	app.Get("/user", handlers.UserGet)
}
