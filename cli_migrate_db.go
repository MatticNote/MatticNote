package main

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	"github.com/urfave/cli/v2"
	"log"
)

func migrateDB(_ *cli.Context) error {
	err := config.LoadConf()
	if err != nil {
		return err
	}

	if err := config.ValidateConfig(); err != nil {
		return err
	}

	err = database.MigrateProcess()
	if err != nil {
		return err
	}

	log.Println("Migrate process successfully.")
	return nil
}
