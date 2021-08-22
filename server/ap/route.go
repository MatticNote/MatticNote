package ap

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(r fiber.Router) {
	r.Get("/inbox", inboxGet)
	r.Post("/inbox", adaptor.HTTPHandlerFunc(inboxPost))
	userGroup := r.Group("/user")
	userGroup.Get("/:uuid", apUserController)
}
