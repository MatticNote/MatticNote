package main

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
	"log"
)

func cliMigrate(c *cli.Context) error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	applied, err := database.DBMigrate(
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.User,
		config.Config.Database.Password,
		config.Config.Database.Name,
		config.Config.Database.SSLMode,
	)
	if err != nil {
		return err
	}

	if applied > 0 {
		log.Panicf("Applied %d migration(s)", applied)
	} else {
		log.Println("No migrations applied")
	}

	return nil
}
