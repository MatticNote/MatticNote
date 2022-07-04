package main

import (
	"fmt"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
)

var mnCli = &cli.App{
	Name:        "MatticNote",
	Description: "Social Networking Service",
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
