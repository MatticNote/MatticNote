package account

import (
	"github.com/MatticNote/MatticNote/internal/types"
	"github.com/gofiber/fiber/v2"
)

func settingRoute(r fiber.Router) {
	r.Get("/core", settingCoreGet)
}

func settingCoreGet(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*types.User)

	return c.Render(
		"account/setting/core",
		fiber.Map{
			"user": user,
		},
	)
}
