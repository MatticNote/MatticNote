package common

import (
	"errors"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/gofiber/fiber/v2"
)

const (
	CSRFFormName   = "csrf_token"
	CSRFContextKey = "csrfToken"
	CSRFCookieName = "_csrf"
)

func CSRFErrorHandler(c *fiber.Ctx, _ error) error {
	return c.Status(fiber.StatusForbidden).Render(
		"error/csrf",
		fiber.Map{},
	)
}

func ValidateCookie(c *fiber.Ctx) error {
	token := c.Cookies(TokenCookieName)
	if token == "" {
		session, err := AccountSession.Get(c)
		if err != nil {
			return err
		}
		session.Set(AccountSessionRedirectTo, c.OriginalURL())
		err = session.Save()
		if err != nil {
			return err
		}
		return c.Redirect("/account/login")
	}

	user, err := account.GetUserFromToken(token)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound),
			errors.Is(err, account.ErrUserGone),
			errors.Is(err, account.ErrInvalidUserToken),
			errors.Is(err, account.ErrUserSuspend):
			return c.Redirect("/account/logout")
		default:
			return err
		}
	}

	c.Locals("currentUser", user)

	return c.Next()
}

func RequireActiveAccount(c *fiber.Ctx) error {
	user, ok := c.Locals("currentUser").(*schemas.User)

	if !ok {
		session, err := AccountSession.Get(c)
		if err != nil {
			return err
		}
		session.Set(AccountSessionRedirectTo, c.OriginalURL())
		err = session.Save()
		if err != nil {
			return err
		}
		return c.Redirect("/account/login")
	}

	if user.DeletedAt.Valid {
		return c.Redirect("/settings/core")
	}

	if !user.Username.Valid {
		return c.Redirect("/settings/core")
	}

	return c.Next()
}

func RequireAdminOrModerator(c *fiber.Ctx) error {
	user, ok := c.Locals("currentUser").(*schemas.User)
	if !ok {
		session, err := AccountSession.Get(c)
		if err != nil {
			return err
		}
		session.Set(AccountSessionRedirectTo, c.OriginalURL())
		err = session.Save()
		if err != nil {
			return err
		}
		return c.Redirect("/account/login")
	}

	if !(user.IsModerator || user.IsAdmin) {
		return fiber.ErrForbidden
	}

	return c.Next()
}
