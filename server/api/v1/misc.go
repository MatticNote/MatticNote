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

func activeAccountRequired(c *fiber.Ctx) error {
	user, ok := c.Locals("currentUser").(*schemas.User)
	if !ok {
		return apiUnauthorized(c)
	}

	if user.DeletedAt.Valid {
		return apiError(c, fiber.StatusForbidden, "ACCOUNT_DELETION_RESERVED", "You are reserving account deletion. If to continue use this account, please cancel account deletion reservation.")
	}

	if !user.Username.Valid {
		return apiError(c, fiber.StatusBadRequest, "USERNAME_REQUIRED", "Please pick a username first.")
	}

	return c.Next()
}
