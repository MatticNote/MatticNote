package misc

import (
	"github.com/gofiber/fiber/v2"
)

func IsAPAcceptHeader(c *fiber.Ctx) bool {
	return c.Accepts("*/*", "application/ld+json", "application/activity+json") != "*/*"
}
