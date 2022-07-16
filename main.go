package main

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
	"time"
)

var mnCli = &cli.App{
	Name:        "MatticNote",
	Description: "Social Networking Service",
	Authors: []*cli.Author{
		{
			Name: "YuzuRyo61",
		},
	},
	Copyright: "(C) 2022 YuzuRyo61",
	Compiled: func() time.Time {
		parse, err := time.Parse("2006/01/02-15:04:05-0700", internal.GetBuildDate())
		if err != nil {
			return time.Now()
		}

		return parse
	}(),
	Version: fmt.Sprintf(
		"%s (Build Date: %s, Go Version: %s)",
		internal.GetSysVersion(),
		internal.GetBuildDate(),
		runtime.Version(),
	),
	EnableBashCompletion: true,
	Commands: []*cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Start server",
			Action:  cliServer,
			Flags: []cli.Flag{
				&cli.PathFlag{
					Name:    "config",
					Aliases: []string{"c"},
					Value:   config.MNConfigDefaultPath,
					Usage:   "MatticNote configuration file",
				},
				&cli.BoolFlag{
					Name:    "logging",
					Aliases: []string{"l"},
					Usage:   "Whether to enable access logging",
				},
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "Migrate database",
			Action:  cliMigrate,
			Flags: []cli.Flag{
				&cli.PathFlag{
					Name:    "config",
					Aliases: []string{"c"},
					Value:   config.MNConfigDefaultPath,
				},
			},
		},
	},
}

func main() {
	if err := mnCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
