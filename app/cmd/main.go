package main

import (
	"fmt"
	"net/http"

	"auth/app/assets"
	"auth/app/cmd/cli"
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
	"github.com/soypat/rebed"
)

func main() {
	database.Connect()
	store.SetStore()

	switch cli.SetCli() {
	case cli.ReturnMigrate:
		database.Migrate(&models.User{})

	case cli.ReturnSuperuser:
		user, pass := config.GetSuperuser()
		hash, err := utils.GetHash(pass)
		if err != nil {
			panic(err)
		}
		u := &models.User{Username: user, Password: string(hash), Admin: true}
		models.CreateUser(u)

	case cli.ReturnStart:
		app := fiber.New(fiber.Config{
			ErrorHandler: utils.ErrorHandler,
			Views:        utils.SetTampletatags(html.NewFileSystem(http.FS(assets.Viewfs), ".html")),
		})

		handlers.SetHtmx(app)
		handlers.SetCsrf(app)

		routes.FreeRoutes(app)
		routes.SuperUserRoutes(app)
		routes.AuthRoutes(app)

		// [static](https://docs.gofiber.io/api/app#static)
		staticConfig := fiber.Static{
			Compress:  true,
			ByteRange: true,
			// Browse:    true,
		}
		rebed.Write(assets.Staticfs, ".")
		app.Static("/static", utils.SetPath(assets.Staticdir), staticConfig)
		app.Static("/media", utils.SetPath(assets.Mediadir), staticConfig)
		app.Get("/metrics", monitor.New(monitor.Config{Title: fmt.Sprintf("%v:%v - metrics", config.IP, config.PORT)}))
		app.Use(logger.New())
		app.Listen(fmt.Sprintf("%v:%v", config.IP, config.PORT))
	}
}

// "github.com/gofiber/fiber/v2/middleware/filesystem"
// [Embed static]
// app.Use("/static", filesystem.New(filesystem.Config{
// 	Root:       http.FS(assets.Staticfs),
// 	PathPrefix: "static",
// 	Browse:     false,
// }))

// "github.com/gofiber/template/django/v3"
// [Django templates]
// engine := django.NewPathForwardingFileSystem(http.FS(assets.Viewfs), "/templates", ".html")
