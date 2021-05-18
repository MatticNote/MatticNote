package account

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
)

func destroySession(c *fiber.Ctx) error {
	if c.Cookies(internal.JWTAuthCookieName, "") != "" {
		c.ClearCookie(internal.JWTAuthCookieName)
		return c.Redirect("/account/login?logout=true")
	} else {
		return c.Redirect("/account/login")
	}
}
