package server

import "github.com/gofiber/fiber/v2"

func ConfigureRoute(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello")
	})
}
