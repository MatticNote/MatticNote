package main

import (
	"embed"
	"fmt"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const (
	DefaultPort = 3000
)

//go:embed static/**
var staticFS embed.FS

//go:embed client/dist/cli/**
var webCliFS embed.FS

var mnAppCli = &cli.App{
	Name:                 "MatticNote",
	Description:          "ActivityPub compatible SNS that aims to be easy for everyone to use",
	Version:              fmt.Sprintf("%s-%s", internal.Version, internal.Revision),
	EnableBashCompletion: true,
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
					Value:       "",
					DefaultText: "",
				},
				&cli.BoolFlag{
					Name:    "skip-migration",
					Usage:   "Start the server without the migration process. Specify when all migrations are applicable.",
					Aliases: []string{"m"},
					EnvVars: []string{"MN_SKIP_MIGRATION"},
				},
				&cli.BoolFlag{
					Name:    "no-worker",
					Usage:   "Launch the web app without launching the worker.",
					Aliases: []string{"w"},
					EnvVars: []string{"MN_NO_WORKER"},
				},
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "Migrate database",
			Action:  migrateDB,
		},
		{
			Name:    "worker",
			Aliases: []string{"w"},
			Usage:   "Start worker",
			Action:  startWorker,
		},
		{
			Name:    "testmail",
			Aliases: []string{"tm"},
			Usage:   "Send a mail for test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "to",
					Required: true,
				},
			},
			Action: testSendMail,
		},
	},
}

func main() {
	if err := mnAppCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
