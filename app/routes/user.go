package routes

import (
	"auth/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func FreeRoutes(app fiber.Router) {
	r := app.Group("/auth")
	r.Get("/login", handlers.LoginPageGet).Name("auth-login")
	r.Post("/login", handlers.LoginPost)
	r.Get("/logout", handlers.LogoutGet).Name("auth-logout")
	r.Get("/health", handlers.Health)
}

func SuperUserRoutes(app fiber.Router) {
	r := app.Group("/auth")
	r.Use(handlers.AuthMiddleware)
	r.Use(handlers.SuperUserMiddleware)
	r.Use(handlers.HtmxMiddleware)
	r.Get("/users", handlers.UserControlGet).Name("auth-users")
	r.Get("/user/:id<int>/edit", handlers.UserEditGet).Name("auth-useredit-get")
	r.Post("/user/:id<int>/edit", handlers.UserEditPost).Name("auth-useredit-post")
	r.Delete("/user/:id<int>/delete", handlers.UserDeletePost).Name("auth-userdelete-post")
	r.Get("/new", handlers.CreateNewUserGet).Name("auth-new")
	r.Post("/new", handlers.CreateNewUserPost)
}

func AuthRoutes(app fiber.Router) {
	r := app.Group("")
	r.Get("/user", handlers.AuthMiddleware, handlers.UserGet).Name("auth-user")
}
