package well_known

import "github.com/gofiber/fiber/v2"

func redirectToChPasswd(c *fiber.Ctx) error {
	return c.Redirect("/account/settings/password", 301)
}
