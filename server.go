package main

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/server"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/urfave/cli/v2"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func cliServer(_ *cli.Context) error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	if err := database.ConnectDB(
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.User,
		config.Config.Database.Password,
		config.Config.Database.Name,
		config.Config.Database.SSLMode,
	); err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
		ServerHeader:          "MatticNote",
		Prefork:               config.Config.Server.Prefork,
		CaseSensitive:         true,
		DisableStartupMessage: true,
		ErrorHandler:          server.ErrorView,
		Views: html.NewFileSystem(func() http.FileSystem {
			dist, err := fs.Sub(template, "template")
			if err != nil {
				panic(err)
			}

			return http.FS(dist)
		}(), ".html"),
	})

	app.Use(recover2.New(recover2.Config{
		EnableStackTrace: true,
	}))

	app.Use("/static/ui", filesystem.New(filesystem.Config{
		Root: func() http.FileSystem {
			webUiDist, err := fs.Sub(webUi, "client/dist/ui")
			if err != nil {
				panic(err)
			}
			return http.FS(webUiDist)
		}(),
		Browse: false,
	}))

	app.Use("/web", filesystem.New(filesystem.Config{
		Root: func() http.FileSystem {
			webCliDist, err := fs.Sub(webCli, "client/dist/client")
			if err != nil {
				panic(err)
			}
			return http.FS(webCliDist)
		}(),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	server.ConfigureRoute(app)

	app.Use(server.NotFoundView)

	database.InitFiberRedisMemory(
		config.Config.Redis.Host,
		config.Config.Redis.Port,
		config.Config.Redis.User,
		config.Config.Redis.Password,
		config.Config.Redis.Database,
	)
	database.InitRedis(
		config.Config.Redis.Host,
		config.Config.Redis.Port,
		config.Config.Redis.User,
		config.Config.Redis.Password,
		config.Config.Redis.Database,
	)

	listen := fmt.Sprintf("%s:%d", config.Config.Server.Host, config.Config.Server.Port)
	if !fiber.IsChild() {
		log.Println(fmt.Sprintf("MatticNote is running at http://%s", listen))
	}

	go func() {
		if err := app.Listen(listen); err != nil {
			panic(err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)

	<-sc
	if !fiber.IsChild() {
		log.Println("MatticNote is shutting down...")
	}

	_ = app.Shutdown()
	_ = database.CloseDB()

	return nil
}
