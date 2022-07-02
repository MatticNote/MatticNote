package account

import (
	"errors"
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
		switch {
		case errors.Is(err, account.ErrInvalidCredentials), errors.Is(err, account.ErrUserGone):
			c.Locals("invalid", true)
			return loginGet(c)
		}
	}

	session, err := account.GenerateUserToken(user.ID, c.IP())
	if err != nil {
		return err
	}

	account.InsertTokenCookie(c, session)

	isEmailVerified, err := account.IsEmailVerified(user.ID)
	if err != nil {
		return err
	}

	if user.DeletedAt.Valid {
		return c.Redirect("/account/settings/core")
	}

	if isEmailVerified {
		if !user.Username.Valid {
			return c.Redirect("/account/register-username")
		} else {
			return c.Redirect("/web")
		}
	} else {
		return c.Redirect("/account/settings/core")
	}
}

func logout(c *fiber.Ctx) error {
	err := account.DeleteUserTokenFromToken(c.Cookies(account.TokenCookieName))
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
