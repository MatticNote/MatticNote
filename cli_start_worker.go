package main

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/worker"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func startWorker(_ *cli.Context) error {
	if err := config.LoadConf(); err != nil {
		return err
	}

	if err := config.ValidateConfig(); err != nil {
		return err
	}

	if err := database.ConnectDB(); err != nil {
		return err
	}

	worker.InitWorker()

	worker.Worker.Start()
	log.Println("MatticNote worker is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)

	_ = <-sc

	log.Println("MatticNote worker is shutting down...")
	worker.Worker.Stop()
	database.DisconnectDB()
	log.Println("MatticNote worker is successful shutdown.")

	return nil
}
