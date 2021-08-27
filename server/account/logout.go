package account

import (
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/gofiber/fiber/v2"
	"time"
)

func destroySession(c *fiber.Ctx) error {
	DestroySessionCookie(c)
	return c.Redirect("/account/login?logout=true", 307)
}

func DestroySessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:    signature.JWTAuthCookieName,
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
	})
}
