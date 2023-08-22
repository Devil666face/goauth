package cmd

import (
	"net/http"

	"app/assets"
	"app/config"
	"app/database"
	"app/middlewares"
	"app/models"
	"app/routes"
	"app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html/v2"
)

func Migrate() error {
	err := database.Migrate(&models.User{})
	if err != nil {
		return err
	}
	return nil
}

func CreateSuperuser() error {
	user, pass := config.GetSuperuser()
	hash, err := utils.GetHash(pass)
	if err != nil {
		return err
	}
	u := &models.User{Username: user, Password: string(hash), Admin: true}
	createrr := models.CreateUser(u)
	if createrr != nil {
		return err
	}
	return nil
}

func StartApp() error {
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
		Views:        utils.SetTampletatags(html.NewFileSystem(http.FS(assets.Viewfs), ".html")),
	})
	app.Use(logger.New())
	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(assets.Staticfs),
		PathPrefix: "static",
		MaxAge:     86400,
	}))
	app.Static("/media", utils.SetPath(assets.Mediadir), fiber.Static{
		Compress:  true,
		ByteRange: true,
		// Browse:    true,
	})
	app.Get("/metrics", monitor.New(monitor.Config{}))

	app.Use(
		middlewares.AllowedHostMiddleware,
		middlewares.HtmxMiddleware,
		middlewares.CsrfMiddleware,
		middlewares.TemplatetagMiddleware,
	)

	routes.FreeRoutes(app)
	routes.SuperUserRoutes(app)
	routes.AuthRoutes(app)

	if config.TLS {
		go func() {
			httpapp := fiber.New(fiber.Config{DisableStartupMessage: true})
			httpapp.Use(middlewares.HttpsRedirectMiddleware)
			err := httpapp.Listen(config.CONNECT_HTTP)
			if err != nil {
				panic(err)
			}
		}()
		err := app.ListenTLS(config.CONNECT_HTTPS, config.TLS_CRT, config.TLS_KEY)
		if err != nil {
			return err
		}
	}
	err := app.Listen(config.CONNECT_HTTP)
	if err != nil {
		return err
	}
	return nil
}

// "github.com/soypat/rebed"
// rebed.Write(assets.Staticfs, ".")
// app.Static("/static", utils.SetPath(assets.Staticdir), staticConfig)

// "github.com/gofiber/template/django/v3"
// [Django templates]
// engine := django.NewPathForwardingFileSystem(http.FS(assets.Viewfs), "/templates", ".html")
