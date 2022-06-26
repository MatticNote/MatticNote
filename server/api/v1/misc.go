package v1

import (
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/gofiber/fiber/v2"
)

func loginRequired(c *fiber.Ctx) error {
	_, ok := c.Locals("currentUser").(*schemas.User)
	if !ok {
		return apiUnauthorized(c)
	}

	return c.Next()
}

func usernameRequired(c *fiber.Ctx) error {
	user, ok := c.Locals("currentUser").(*schemas.User)
	if !ok {
		return apiUnauthorized(c)
	}

	if !user.Username.Valid {
		return apiError(c, fiber.StatusBadRequest, "USERNAME_REQUIRED", "Please pick a username first.")
	}

	return c.Next()
}
