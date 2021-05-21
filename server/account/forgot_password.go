package account

import (
	"context"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"net/http"
)

type issueForgotPasswordFormStruct struct {
	Email string `validate:"required,email"`
}

type forgotPasswordFormStruct struct {
	Password string `validate:"required,min=8"`
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
	formData := new(issueForgotPasswordFormStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(*formData); errs != nil {
		return forgotPasswordView(c, errs...)
	}

	if err := internal.IssueForgotPassword(formData.Email); err != nil {
		return err
	}

	return c.Redirect("/account/login?forgot_sent=true")
}

func forgotPasswordResetGet(c *fiber.Ctx) error {
	_, err := internal.ValidateForgotPasswordToken(c.Params("token"))
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(http.StatusBadRequest).Render(
				"account/forgot_reset_error",
				fiber.Map{},
				"_layout/account",
			)
		} else {
			return err
		}
	}
	return forgotPasswordResetView(c)
}

func forgotPasswordResetView(c *fiber.Ctx, errors ...string) error {
	return c.Status(http.StatusOK).Render(
		"account/forgot_reset",
		fiber.Map{
			"CSRFFormName": csrfFormName,
			"CSRFToken":    c.Context().UserValue(csrfContextKey).(string),
			"errors":       errors,
		},
		"_layout/account",
	)
}

func forgotPasswordResetPost(c *fiber.Ctx) error {
	targetUuid, err := internal.ValidateForgotPasswordToken(c.Params("token"))
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(http.StatusBadRequest).Render(
				"account/forgot_reset_error",
				fiber.Map{},
				"_layout/account",
			)
		} else {
			return err
		}
	}

	formData := new(forgotPasswordFormStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(*formData); errs != nil {
		return forgotPasswordResetView(c, errs...)
	}

	if err := internal.ChangeUserPassword(targetUuid, formData.Password); err != nil {
		return err
	}

	_, err = database.DBPool.Exec(
		context.Background(),
		"delete from user_reset_password where target = $1;",
		targetUuid,
	)
	if err != nil {
		return err
	}

	return c.Redirect("/account/login?password_reset=true")
}
