package account

import (
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/gofiber/fiber/v2"
)

func settingRoute(r fiber.Router) {
	r.Get("/core", settingCoreGet)
}

func settingCoreGet(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)

	isEmailVerified, err := account.IsEmailVerified(user.ID)
	if err != nil {
		return err
	}
	email, err := account.GetUserEmail(user.ID)
	if err != nil {
		return err
	}

	return c.Render(
		"account/setting/core",
		fiber.Map{
			"user":                user,
			"isUserEmailVerified": isEmailVerified,
			"userEmail":           email,
		},
	)
}
