package settings

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(r fiber.Router) {
	r.Use(internal.RegisterFiberJWT("cookie", true))

	r.Get("/profile", editProfileGet)
	r.Post("/profile", editProfilePost)
}
