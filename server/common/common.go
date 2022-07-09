package common

import (
	"errors"
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
	token := c.Cookies(account.TokenCookieName)
	if token == "" {
		return c.Redirect("/account/login")
	}

	user, err := account.GetUserFromToken(token)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound),
			errors.Is(err, account.ErrUserGone),
			errors.Is(err, account.ErrInvalidUserToken):
			return c.Redirect("/account/logout")
		default:
			return err
		}
	}

	c.Locals("currentUser", user)

	return c.Next()
}
