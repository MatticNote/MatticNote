package main

import (
	"embed"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/mn_template"
	"github.com/MatticNote/MatticNote/server"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	fr "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/django"
	"github.com/urfave/cli/v2"
	"io/fs"
	"log"
	"net/http"
	"os"
)

const (
	DefaultPort = 3000
	DefaultAddr = "127.0.0.1"
)

//go:embed static/**
var staticFS embed.FS

//go:embed client/dist/cli/**
var webCliFS embed.FS

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

	if !fiber.IsChild() {
		if !c.Bool("skip-migration") {
			err := database.MigrateProcess()
			if err != nil {
				return err
			}
		}

		err = internal.GenerateJWTSignKey(false)
		if err != nil {
			return err
		}
	}

	if err := database.ConnectDB(); err != nil {
		return err
	}
	defer database.DisconnectDB()

	err = internal.LoadJWTSignKey()
	if err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
		Prefork:               true,
		ServerHeader:          "MatticNote",
		CaseSensitive:         true,
		Views:                 django.NewFileSystem(http.FS(mn_template.Templates), ".django"),
		ErrorHandler:          server.ErrorView,
		DisableStartupMessage: true,
	})

	app.Use(fr.New(fr.Config{
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

	app.Use("/web", internal.RegisterFiberJWT("cookie"), filesystem.New(filesystem.Config{
		Root: func() http.FileSystem {
			webCliFSDist, err := fs.Sub(webCliFS, "client/dist/cli")
			if err != nil {
				panic(err)
			}
			return http.FS(webCliFSDist)
		}(),
		Browse: false,
	}))

	app.Use(server.NotFoundView)

	if addr == DefaultAddr {
		addr = config.Config.Server.ListenAddress
	}

	if addrPort == DefaultPort {
		addrPort = uint(config.Config.Server.ListenPort)
	}

	listen := fmt.Sprintf("%s:%d", addr, addrPort)
	if !fiber.IsChild() {
		log.Println(fmt.Sprintf("MatticNote is running at http://%s", listen))
	}

	if err := app.Listen(listen); err != nil {
		panic(err)
	}

	return nil
}

func main() {
	if err := mnAppCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
