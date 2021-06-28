package nodeinfo

import "github.com/gofiber/fiber/v2"

func ConfigureRoute(r fiber.Router) {
	r.Get("/2.1", nodeinfoV21)
}
