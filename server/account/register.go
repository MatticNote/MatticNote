package account

import (
	ia "github.com/MatticNote/MatticNote/internal/account"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type RegisterForm struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
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
	form := new(RegisterForm)
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
