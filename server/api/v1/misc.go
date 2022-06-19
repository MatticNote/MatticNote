package v1

import (
	"github.com/MatticNote/MatticNote/internal/types"
	"github.com/gofiber/fiber/v2"
)

func loginRequired(c *fiber.Ctx) error {
	_, ok := c.Locals("currentUser").(*types.User)
	if !ok {
		return apiUnauthorized(c)
	}

	return c.Next()
}
