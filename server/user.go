package server

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/gofiber/fiber/v2"
)

func getUserProfile(c *fiber.Ctx) error {
	user, err := account.GetUserByUsername(c.Params("username"))
	if err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound):
			return fiber.ErrNotFound
		default:
			return err
		}
	}

	return c.Render(
		"user-profile",
		fiber.Map{
			"user": user,
		},
	)
}
