package server

import (
	"github.com/MatticNote/MatticNote/server/api"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hi")
	})
	api.ConfigureRoute(app.Group("/api"))
}
