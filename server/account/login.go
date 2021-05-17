package account

import (
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type loginUserStruct struct {
	Login    string `validate:"required"`
	Password string `validate:"required"`
}

func loginUserGet(c *fiber.Ctx) error {
	return loginUserView(c)
}

func loginUserView(c *fiber.Ctx, errors ...string) error {
	return c.Status(http.StatusOK).Render(
		"account/login",
		fiber.Map{
			"CSRFFormName": csrfFormName,
			"CSRFToken":    c.Context().UserValue(csrfContextKey).(string),
			"errors":       errors,
		},
		"_layout/account",
	)
}

func loginPost(c *fiber.Ctx) error {
	formData := new(loginUserStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(*formData); errs != nil {
		return loginUserView(c, errs...)
	}

	return c.SendString("OK")
}
