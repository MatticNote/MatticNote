package account

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type forgotPasswordFormStruct struct {
	Email string `validate:"required,email"`
}

func forgotPasswordGet(c *fiber.Ctx) error {
	if c.Cookies(internal.JWTAuthCookieName, "") != "" {
		return c.Redirect("/web/", 307)
	}

	return forgotPasswordView(c)
}

func forgotPasswordView(c *fiber.Ctx, errors ...string) error {
	return c.Status(http.StatusOK).Render(
		"account/forgot",
		fiber.Map{
			"CSRFFormName": csrfFormName,
			"CSRFToken":    c.Context().UserValue(csrfContextKey).(string),
			"errors":       errors,
		},
		"_layout/account",
	)
}

func forgotPasswordPost(c *fiber.Ctx) error {
	formData := new(forgotPasswordFormStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(*formData); errs != nil {
		return forgotPasswordView(c, errs...)
	}

	// TODO: パスワードを忘れた時のトークンとかを発行する関数を作る

	return c.Redirect("/account/login?forgot_sent=true")
}
