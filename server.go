package main

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/server"
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func cliServer(c *cli.Context) error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
		ServerHeader:          "MatticNote",
		Prefork:               config.Config.Server.Prefork,
		CaseSensitive:         true,
		DisableStartupMessage: true,
		ErrorHandler:          server.ErrorView,
	})

	app.Use(recover2.New(recover2.Config{
		EnableStackTrace: true,
	}))

	listen := fmt.Sprintf("%s:%d", config.Config.Server.Host, config.Config.Server.Port)
	if !fiber.IsChild() {
		log.Println(fmt.Sprintf("MatticNote is running at http://%s", listen))
	}

	server.ConfigureRoute(app)
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

	return nil
}
