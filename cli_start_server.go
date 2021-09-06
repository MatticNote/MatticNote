package main

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal/oauth"
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/MatticNote/MatticNote/mn_template"
	"github.com/MatticNote/MatticNote/server"
	"github.com/MatticNote/MatticNote/worker"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	gfRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/django"
	"github.com/urfave/cli/v2"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func startServer(c *cli.Context) error {
	var (
		addr     = c.String("address")
		addrPort = c.Uint("port")
	)

	err := config.LoadConf()
	if err != nil {
		return err
	}

	if !fiber.IsChild() {
		if err := config.ValidateConfig(); err != nil {
			return err
		}

		if !c.Bool("skip-migration") {
			err := database.MigrateProcess()
			if err != nil {
				return err
			}
		}

		err = signature.GenerateJWTSignKey(false)
		if err != nil {
			return err
		}
	}

	if err := database.ConnectDB(); err != nil {
		return err
	}

	err = signature.LoadJWTSignKey()
	if err != nil {
		return err
	}

	if !fiber.IsChild() {
		if err := signature.VerifyRSASign(); err != nil {
			return err
		}
	}

	app := fiber.New(fiber.Config{
		Prefork:               config.Config.Server.Prefork,
		ServerHeader:          "MatticNote",
		CaseSensitive:         true,
		Views:                 django.NewFileSystem(http.FS(mn_template.Templates), ".django"),
		ErrorHandler:          server.ErrorView,
		DisableStartupMessage: true,
	})

	app.Use(gfRecover.New(gfRecover.Config{
		EnableStackTrace: true,
	}))

	server.ConfigureRoute(app)

	app.Use("/static", filesystem.New(filesystem.Config{
		Root: func() http.FileSystem {
			staticFSDist, err := fs.Sub(staticFS, "static")
			if err != nil {
				panic(err)
			}
			return http.FS(staticFSDist)
		}(),
		Browse: false,
	}))

	app.Use("/web", signature.RegisterFiberJWT("cookie", true), filesystem.New(filesystem.Config{
		Root: func() http.FileSystem {
			webCliFSDist, err := fs.Sub(webCliFS, "client/dist/cli")
			if err != nil {
				panic(err)
			}
			return http.FS(webCliFSDist)
		}(),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	app.Use(server.NotFoundView)

	if addr == "" {
		addr = config.Config.Server.ListenAddress
	}

	if addrPort == DefaultPort {
		addrPort = uint(config.Config.Server.ListenPort)
	}

	worker.InitEnqueue()
	if !c.Bool("no-worker") {
		worker.InitWorker()
	}

	oauth.InitOAuth()

	if !fiber.IsChild() {
		listenAddr := addr
		if addr == "" {
			listenAddr = "0.0.0.0"
		}
		log.Println(fmt.Sprintf("MatticNote is running at http://%s:%d", listenAddr, addrPort))
	}

	listen := fmt.Sprintf("%s:%d", addr, addrPort)
	go func() {
		if !c.Bool("no-worker") {
			if !fiber.IsChild() {
				worker.Worker.Start()
				log.Println("MatticNote worker is running")
			}
		}
		if err := app.Listen(listen); err != nil {
			panic(err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)

	_ = <-sc
	if !fiber.IsChild() {
		log.Println("MatticNote is shutting down...")
	}

	_ = app.Shutdown()
	if !c.Bool("no-worker") {
		if !fiber.IsChild() {
			worker.Worker.Stop()
			log.Println("MatticNote worker is successful shutdown.")
		}
	}
	database.DisconnectDB()

	if !fiber.IsChild() {
		fmt.Println("MatticNote is successful shutdown.")
	}

	return nil
}
