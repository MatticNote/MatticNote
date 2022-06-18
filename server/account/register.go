package account

import (
	ia "github.com/MatticNote/MatticNote/internal/account"
	"github.com/MatticNote/MatticNote/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type registerForm struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

type registerUsernameForm struct {
	Username string `validate:"required"`
}

func registerGet(c *fiber.Ctx) error {
	invalid, ok := c.Locals("invalid").(bool)
	if !ok {
		invalid = false
	}

	return c.Render("account/register", fiber.Map{
		"invalid":    invalid,
		"title":      "Register",
		"csrf_name":  csrfFormName,
		"csrf_token": c.Locals(csrfContextKey),
	}, "account/_common")
}

func registerPost(c *fiber.Ctx) error {
	form := new(registerForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	err := validator.New().Struct(*form)
	if err != nil {
		c.Locals("invalid", true)
		return registerGet(c)
	}

	account, err := ia.RegisterLocalAccount(
		form.Email,
		form.Password,
		false,
	)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok {
			switch dbErr.Code.Name() {
			case "unique_violation":
				c.Locals("invalid", true)
				return registerGet(c)
			default:
				return err
			}
		} else {
			return err
		}
	}
	return c.SendString(account.ID.String())
}

func verifyEmailToken(c *fiber.Ctx) error {
	token := c.Params("token")

	err := ia.VerifyEmailToken(token)
	if err != nil {
		if err == ia.ErrInvalidToken {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid or expired token")
		} else {
			return err
		}
	}

	return c.Redirect("/account/register-username")
}

func registerUsernameGet(c *fiber.Ctx) error {
	currentUser, ok := c.Locals("currentUser").(*types.User)
	if !ok {
		return c.Redirect("/account/logout")
	}

	if currentUser.Username.Valid {
		return c.Redirect("/")
	}

	if !currentUser.EmailVerified.Valid || !currentUser.EmailVerified.Bool {
		return c.SendStatus(fiber.StatusForbidden)
	}

	invalid, ok := c.Locals("invalid").(bool)
	if !ok {
		invalid = false
	}

	return c.Render("account/register-username", fiber.Map{
		"invalid":    invalid,
		"title":      "Register username",
		"csrf_name":  csrfFormName,
		"csrf_token": c.Locals(csrfContextKey),
	}, "account/_common")
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

	currentUser, ok := c.Locals("currentUser").(*types.User)
	if !ok {
		return c.Redirect("/account/logout")
	}

	if currentUser.Username.Valid {
		return c.SendStatus(fiber.StatusForbidden)
	}

	err = ia.ChooseUsername(currentUser.ID, form.Username)
	if err != nil {
		switch err {
		case ia.ErrUsernameAlreadyTaken:
			c.Locals("invalid", true)
			return registerUsernameGet(c)
		default:
			return err
		}
	}

	return c.Redirect("/")
}
