package misc

import (
	"github.com/go-playground/validator"
	"regexp"
)

func RegisterCommonValidator(v *validator.Validate) {
	_ = v.RegisterValidation("username", validatorUsername)
}

func validatorUsername(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,32}$`).Match([]byte(fl.Field().String()))
}
