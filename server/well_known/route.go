package well_known

import "github.com/gofiber/fiber/v2"

func ConfigureRoute(r fiber.Router) {
	r.Get("/change-password", redirectChPasswd)
}
