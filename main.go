package main

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/mn_template"
	"github.com/MatticNote/MatticNote/server"
	"github.com/gofiber/fiber/v2"
	fr "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/pug"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
)

const (
	DefaultPort = 3000
	DefaultAddr = "127.0.0.1"
)

var mnAppCli = &cli.App{
	Name:        "MatticNote",
	Description: "ActivityPub compatible SNS that aims to be easy for everyone to use",
	Commands: []*cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Start server",
			Action:  startServer,
			Flags: []cli.Flag{
				&cli.UintFlag{
					Name:        "port",
					Usage:       "Specifies the port number for listening to the server",
					Aliases:     []string{"p"},
					EnvVars:     []string{"MN_PORT"},
					Value:       DefaultPort,
					DefaultText: "3000",
				},
				&cli.StringFlag{
					Name:        "address",
					Usage:       "Specified the address for listening to the server",
					Aliases:     []string{"a"},
					EnvVars:     []string{"MN_ADDR"},
					Value:       DefaultAddr,
					DefaultText: "127.0.0.1",
				},
				&cli.BoolFlag{
					Name:    "skip-migration",
					Usage:   "Start the server without the migration process. Specify when all migrations are applicable.",
					Aliases: []string{"m"},
					EnvVars: []string{"MN_SKIP_MIGRATION"},
				},
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "Migrate database",
			Action:  migrateDB,
		},
	},
}

func migrateDB(_ *cli.Context) error {
	err := config.LoadConf()
	if err != nil {
		return err
	}

	err = database.MigrateProcess()
	if err != nil {
		return err
	}

	log.Println("Migrate process successfully.")
	return nil
}

func startServer(c *cli.Context) error {
	var (
		addr     = c.String("address")
		addrPort = c.Uint("port")
	)

	err := config.LoadConf()
	if err != nil {
		return err
	}

	if !c.Bool("skip-migration") {
		err := database.MigrateProcess()
		if err != nil {
			return err
		}
	}

	if err := database.ConnectDB(); err != nil {
		return err
	}
	defer database.DisconnectDB()

	app := fiber.New(fiber.Config{
		Prefork:       true,
		ServerHeader:  "MatticNote",
		CaseSensitive: true,
		ErrorHandler:  server.ErrorView,
		Views:         pug.NewFileSystem(http.FS(mn_template.Templates), ".pug"),
	})

	app.Use(fr.New(fr.Config{
		EnableStackTrace: true,
	}))

	server.ConfigureRoute(app)
	app.Use(server.NotFoundView)

	if err := app.Listen(fmt.Sprintf("%s:%d", addr, addrPort)); err != nil {
		panic(err)
	}

	return nil
}

func main() {
	if err := mnAppCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
