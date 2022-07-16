package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

const MNConfigDefaultPath = "matticnote.toml"

var Config *MNConfig

func LoadConfig(filename string) error {
	file, err := ioutil.ReadFile(filename)
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
