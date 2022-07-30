package common

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

const (
	TokenCookieName = "mn_token"
)

func InsertTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     TokenCookieName,
		Value:    token,
		Path:     "/",
		Secure:   false,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
}

func DestroyTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:    TokenCookieName,
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
	})
}
