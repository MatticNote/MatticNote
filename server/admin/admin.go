package admin

import "github.com/gofiber/fiber/v2"

func ConfigureRoute(r fiber.Router) {
	r.Get("/", adminDashboard)
}

func adminDashboard(c *fiber.Ctx) error {
	return c.Render(
		"admin/index",
		fiber.Map{},
		"admin/_layout",
	)
}
