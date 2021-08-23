package config

import (
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/go-playground/validator"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

var Config *MNConfig

func LoadConf() error {
	file, err := ioutil.ReadFile("matticnote.toml")
	if err != nil {
		return err
	}

	err = toml.Unmarshal(file, &Config)
	if err != nil {
		return err
	}

	return nil
}

func ValidateConfig() error {
	validate := validator.New()
	misc.RegisterCommonValidator(validate)
	err := validate.Struct(Config)
	if err != nil {
		var returnErrStr = "There is a problem with the settings: "
		for _, err := range err.(validator.ValidationErrors) {
			returnErrStr += fmt.Sprintf("%s, ", err.StructNamespace())
		}
		return errors.New(returnErrStr)
	}

	if (Config.Server.RecaptchaSiteKey == "" && Config.Server.RecaptchaSecretKey != "") || (Config.Server.RecaptchaSiteKey != "" && Config.Server.RecaptchaSecretKey == "") {
		return errors.New("missing reCAPTCHA site key or secret key")
	}

	return nil
}
