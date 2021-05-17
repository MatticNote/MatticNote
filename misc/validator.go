package misc

import (
	"fmt"
	"github.com/go-playground/validator"
	"regexp"
)

func RegisterCommonValidator(v *validator.Validate) {
	_ = v.RegisterValidation("username", validatorUsername)
}

func validatorUsername(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,32}$`).Match([]byte(fl.Field().String()))
}

func ValidateForm(data interface{}) []string {
	var errors []string
	validate := validator.New()
	RegisterCommonValidator(validate)
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Invalid form: %s", err.Field()))
		}
	}

	return errors
}
