package ap

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(r fiber.Router) {
	r.Get("/inbox", inboxGet)
	r.Post("/inbox", adaptor.HTTPHandlerFunc(inboxPost))
	ug := r.Group("/user/:uuid", apUserHandler)
	ug.Get("/", apUserController)
	ug.Get("/inbox", inboxGet)
	ug.Post("/inbox", adaptor.HTTPHandlerFunc(inboxPost))
}
