package account

import (
	"github.com/gofiber/fiber/v2"
)

func destroySession(c *fiber.Ctx) error {
	c.ClearCookie()
	return c.Redirect("/account/login?logout=true")
}
