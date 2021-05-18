package account

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
	"time"
)

func destroySession(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:    internal.JWTAuthCookieName,
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
	})
	return c.Redirect("/account/login?logout=true", 307)
}
