package ap

import "github.com/gofiber/fiber/v2"

func ConfigureRoute(r fiber.Router) {
	r.Post("/inbox", inboxPost)
}
