package main

import (
	"fmt"
	"net/http"

	"auth/app/assets"
	"auth/app/config"
	"auth/app/database"
	"auth/app/handlers"
	"auth/app/models"
	"auth/app/routes"
	"auth/app/store"
	"auth/app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html/v2"
)

func main() {
	database.Connect()
	database.Migrate(&models.User{})

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
		Views:        html.NewFileSystem(http.FS(assets.Viewfs), ".html"),
	})

	store.SetStore()

	app.Use(logger.New())

	handlers.SetHtmx(app)
	handlers.SetCsrf(app)

	routes.FreeRoutes(app)
	routes.SuperUserRoutes(app)
	routes.AuthRoutes(app)

	app.Get("/metrics", monitor.New(monitor.Config{Title: fmt.Sprintf("%v:%v - metrics", config.IP, config.PORT)}))
	app.Listen(fmt.Sprintf("%v:%v", config.IP, config.PORT))
}
