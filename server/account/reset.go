package account

import (
	"errors"
	ia "github.com/MatticNote/MatticNote/internal/account"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	resetPasswordNewStruct struct {
		Email string `validate:"email,required"`
	}

	resetPasswordFormStruct struct {
		Password string `validate:"required,min=8"`
	}
)

func resetPasswordNewGet(c *fiber.Ctx) error {
	if c.Cookies(ia.TokenCookieName) != "" {
		return c.Redirect("/web")
	}

	return c.Render(
		"account/reset-password",
		fiber.Map{
			"csrfName":  csrfFormName,
			"csrfToken": c.Locals(csrfContextKey),
		},
		"account/_layout",
	)
}

func resetPasswordNewPost(c *fiber.Ctx) error {
	form := new(resetPasswordNewStruct)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	err := validator.New().Struct(*form)
	if err != nil {
		c.Locals("invalid", true)
		return resetPasswordNewGet(c)
	}

	err = ia.NewPasswordResetToken(form.Email)
	if err != nil {
		return err
	}

	return c.Render(
		"account/reset-password_post",
		fiber.Map{},
		"account/_layout",
	)
}

func resetPasswordFormGet(c *fiber.Ctx) error {
	token := c.Params("token")

	if !ia.IsValidPasswordResetToken(token) {
		return c.Status(fiber.StatusBadRequest).Render(
			"account/reset-password_invalid",
			fiber.Map{},
			"account/_layout",
		)
	}

	return c.Render(
		"account/reset-password_form",
		fiber.Map{
			"csrfName":  csrfFormName,
			"csrfToken": c.Locals(csrfContextKey),
		},
		"account/_layout",
	)
}

func resetPasswordFormPost(c *fiber.Ctx) error {
	token := c.Params("token")

	form := new(resetPasswordFormStruct)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	err := validator.New().Struct(*form)
	if err != nil {
		c.Locals("invalid", true)
		return resetPasswordFormGet(c)
	}

	userId, err := ia.PopPasswordResetToken(token)
	if err != nil {
		if errors.Is(err, ia.ErrInvalidPasswordResetToken) {
			return c.Status(fiber.StatusBadRequest).Render(
				"account/reset-password_invalid",
				fiber.Map{},
				"account/_layout",
			)
		} else {
			return err
		}
	}

	err = ia.UpdateUserPassword(*userId, form.Password)
	if err != nil {
		return err
	}

	return c.Render(
		"account/reset-password_complete",
		fiber.Map{},
		"account/_layout",
	)
}
