package admin

import (
	"fmt"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func ConfigureRoute(r fiber.Router) {
	r.Get("/", adminDashboard)
	r.Get("/metrics", monitor.New(monitor.Config{
		Title: fmt.Sprintf("MatticNote %s", internal.GetSysVersion()),
	}))
}

func adminDashboard(c *fiber.Ctx) error {
	return c.Render(
		"admin/index",
		fiber.Map{},
		"admin/_layout",
	)
}
