package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var mnCli = &cli.App{
	Name:                 "MatticNote",
	Description:          "Social Networking Service",
	Version:              fmt.Sprintf("0.0.0"),
	EnableBashCompletion: true,
	Commands: []*cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Start server",
			Action:  cliServer,
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "Migrate database",
			Action:  cliMigrate,
		},
	},
}

func main() {
	if err := mnCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
