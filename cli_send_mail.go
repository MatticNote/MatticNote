package main

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal/mail"
	"github.com/urfave/cli/v2"
	"log"
)

func testSendMail(c *cli.Context) error {
	err := config.LoadConf()
	if err != nil {
		return err
	}

	if err := config.ValidateConfig(); err != nil {
		return err
	}

	err = mail.SendMail(
		c.String("to"),
		"MatticNote Test mail / MatticNote テストメール",
		"text/plain",
		"If you can see this mail, configuration is correct!\n"+
			"このメッセージが見えている場合、設定は正しいです！",
	)
	if err == nil {
		log.Println("Test mail was sent!")
	}
	return err
}
