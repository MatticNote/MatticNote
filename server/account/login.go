package account

import (
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"
)

type loginForm struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

func loginGet(c *fiber.Ctx) error {
	if c.Cookies(account.TokenCookieName) != "" {
		return c.Redirect("/web/")
	}

	return c.Render("account/login", fiber.Map{
		"csrf_name":  csrfFormName,
		"csrf_token": c.Locals(csrfContextKey),
	})
}

func loginPost(c *fiber.Ctx) error {
	form := new(loginForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	err := validator.New().Struct(*form)
	if err != nil {
		c.Locals("invalid", true)
		return loginGet(c)
	}

	user, err := account.AuthenticateUser(form.Email, form.Password)
	if err != nil {
		if err == account.ErrInvalidCredentials {
			return c.SendStatus(fiber.StatusBadRequest)
		} else {
			return err
		}
	}

	session, err := account.GenerateUserToken(user.ID, c.IP())
	if err != nil {
		return err
	}

	account.InsertTokenCookie(c, session)

	if !user.Username.Valid && (user.EmailVerified.Valid && user.EmailVerified.Bool) {
		return c.Redirect("/account/register-username")
	} else if !user.EmailVerified.Bool {
		return c.Redirect("/account/settings/core")
	} else {
		return c.Redirect("/web/")
	}
}

func logout(c *fiber.Ctx) error {
	err := account.DestroyUserToken(c.Cookies(account.TokenCookieName))
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    account.TokenCookieName,
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
	})
	return c.Redirect("/account/login", fiber.StatusTemporaryRedirect)
}
