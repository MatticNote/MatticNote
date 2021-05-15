package account

import (
	"github.com/MatticNote/MatticNote/misc"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type registerUserStruct struct {
	Username string `validate:"required,username"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
	// TODO: CAPTCHAなどの対策用のフォーム内容も含める
}

func registerUserGet(c *fiber.Ctx) error {
	return registerUserView(c)
}

func registerUserView(c *fiber.Ctx, errors ...string) error {
	return c.Status(http.StatusOK).Render(
		"register",
		fiber.Map{
			"CSRFFormName": csrfFormName,
			"CSRFToken":    c.Context().UserValue(csrfContextKey).(string),
			"errors":       errors,
		},
		"layout/account",
	)
}

func registerPost(c *fiber.Ctx) error {
	formData := new(registerUserStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := validateForm(*formData); errs != nil {
		return registerUserView(c, errs...)
	}

	// TODO: アカウント作成関数とか

	return c.Status(200).SendString("OK")
	//return c.Redirect("/account/login?created")
}

func validateForm(data registerUserStruct) []string {
	var errors []string
	validate := validator.New()
	misc.RegisterCommonValidator(validate)
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field())
		}
	}

	return errors
}
