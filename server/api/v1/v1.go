package v1

import "github.com/gofiber/fiber/v2"

func ConfigureRoute(r fiber.Router) {
	r.Use(func(c *fiber.Ctx) error {
		// TODO: User information
		return c.Next()
	})

	userApiRoute(r.Group("/users"))
}
