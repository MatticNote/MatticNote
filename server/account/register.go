package account

import (
	"errors"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/database/schemas"
	ia "github.com/MatticNote/MatticNote/internal/account"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type registerForm struct {
	Email      string `validate:"required,email"`
	Password   string `validate:"required,min=8"`
	InviteCode string `json:"invite_code"`
}

type registerUsernameForm struct {
	Username string `validate:"required"`
}

func registerGet(c *fiber.Ctx) error {
	if config.Config.System.RegistrationMode == 0 {
		return fiber.ErrForbidden
	}

	var requiredInviteCode = false
	if config.Config.System.RegistrationMode == 1 {
		requiredInviteCode = true
	}

	invalid, ok := c.Locals("invalid").(bool)
	if !ok {
		invalid = false
	}

	if c.Cookies(ia.TokenCookieName) != "" {
		return c.Redirect("/web")
	}

	return c.Render("account/register", fiber.Map{
		"invalid":            invalid,
		"title":              "Register",
		"csrfName":           csrfFormName,
		"csrfToken":          c.Locals(csrfContextKey),
		"requiredInviteCode": requiredInviteCode,
	}, "account/_layout")
}

func registerPost(c *fiber.Ctx) error {
	if config.Config.System.RegistrationMode == 0 {
		return fiber.ErrForbidden
	}

	form := new(registerForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	err := validator.New().Struct(*form)
	if err != nil {
		c.Locals("invalid", true)
		return registerGet(c)
	}

	if config.Config.System.RegistrationMode == 1 {
		// TODO: Invite token use method
	}

	_, err = ia.RegisterLocalAccount(
		form.Email,
		form.Password,
		false,
	)
	if err != nil {
		if errors.Is(err, ia.ErrEmailExists) {
			c.Locals("invalid", true)
			return registerGet(c)
		} else {
			return err
		}
	}

	return c.Render("account/register_post", fiber.Map{}, "account/_layout")
}

func verifyEmailToken(c *fiber.Ctx) error {
	token := c.Params("token")

	err := ia.VerifyEmailToken(token)
	if err != nil {
		if errors.Is(err, ia.ErrInvalidToken) {
			return c.Status(fiber.StatusBadRequest).Render(
				"account/invalid-email-token",
				fiber.Map{},
				"account/_layout",
			)
		} else {
			return err
		}
	}

	currentUser, ok := c.Locals("currentUser").(*schemas.User)
	if !ok {
		return c.Redirect("/account/login")
	}

	if !currentUser.Username.Valid {
		return c.Redirect("/account/register-username")
	} else {
		return c.Redirect("/web")
	}
}

func registerUsernameGet(c *fiber.Ctx) error {
	currentUser, ok := c.Locals("currentUser").(*schemas.User)
	if !ok {
		return c.Redirect("/account/logout")
	}

	if currentUser.Username.Valid {
		return c.Redirect("/web/")
	}

	verified, err := ia.IsEmailVerified(currentUser.ID)
	if err != nil {
		return err
	}
	if !verified || currentUser.DeletedAt.Valid {
		return c.SendStatus(fiber.StatusForbidden)
	}

	invalid, ok := c.Locals("invalid").(bool)
	if !ok {
		invalid = false
	}

	return c.Render("account/register-username", fiber.Map{
		"invalid":   invalid,
		"title":     "Register username",
		"csrfName":  csrfFormName,
		"csrfToken": c.Locals(csrfContextKey),
	}, "account/_layout")
}

func registerUsernamePost(c *fiber.Ctx) error {
	form := new(registerUsernameForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	err := validator.New().Struct(*form)
	if err != nil {
		c.Locals("invalid", true)
		return registerUsernameGet(c)
	}

	currentUser, ok := c.Locals("currentUser").(*schemas.User)
	if !ok {
		return c.Redirect("/account/logout")
	}

	verified, err := ia.IsEmailVerified(currentUser.ID)
	if err != nil {
		return err
	}
	if !verified || currentUser.DeletedAt.Valid {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if currentUser.Username.Valid {
		return c.SendStatus(fiber.StatusForbidden)
	}

	err = ia.ChooseUsername(currentUser.ID, form.Username)
	if err != nil {
		switch {
		case errors.Is(err, ia.ErrUsernameAlreadyTaken):
			c.Locals("invalid", true)
			return registerUsernameGet(c)
		case errors.Is(err, ia.ErrInvalidUsernameFormat):
			c.Locals("invalid", true)
			return registerUsernameGet(c)
		default:
			return err
		}
	}

	return c.Redirect("/web/")
}
