package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

var Config *MNConfig

func LoadConfig() error {
	file, err := ioutil.ReadFile("matticnote.toml")
	if err != nil {
		return err
	}

	var cfg MNConfig

	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		return err
	}

	err = validator.New().Struct(cfg)
	if err != nil {
		return err
	}

	Config = &cfg
	return nil
}
