package account

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/MatticNote/MatticNote/server/common"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type loginForm struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func loginGet(c *fiber.Ctx) error {
	if c.Cookies(common.TokenCookieName) != "" {
		return c.Redirect("/web/")
	}

	return c.Render("account/login", fiber.Map{
		"csrfName":  common.CSRFFormName,
		"csrfToken": c.Locals(common.CSRFContextKey),
		"invalid":   c.Locals("invalid"),
		"suspend":   c.Locals("suspend"),
	}, "account/_layout")
}

func loginPost(c *fiber.Ctx) error {
	form := new(loginForm)
	if err := c.BodyParser(form); err != nil {
		return fiber.ErrUnprocessableEntity
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
		case errors.Is(err, account.ErrUserSuspend):
			c.Locals("suspend", true)
			return loginGet(c)
		}
	}

	session, err := account.GenerateUserToken(user.ID, c.IP())
	if err != nil {
		return err
	}

	common.InsertTokenCookie(c, session)

	isEmailVerified, err := account.IsEmailVerified(user.ID)
	if err != nil {
		return err
	}

	if user.DeletedAt.Valid {
		return c.Redirect("/settings/core")
	}

	if isEmailVerified {
		if !user.Username.Valid {
			return c.Redirect("/account/register-username")
		} else {
			session, err := common.AccountSession.Get(c)
			if err != nil {
				return err
			}
			redirectTo, ok := session.Get(common.AccountSessionRedirectTo).(string)
			if !ok {
				return c.Redirect("/web")
			} else {
				err := session.Destroy()
				if err != nil {
					return err
				}
				return c.Redirect(redirectTo)
			}
		}
	} else {
		return c.Redirect("/settings/core")
	}
}

func logout(c *fiber.Ctx) error {
	err := account.DeleteUserTokenFromToken(c.Cookies(common.TokenCookieName))
	if err != nil {
		return err
	}

	common.DestroyTokenCookie(c)
	return c.Redirect("/")
}
