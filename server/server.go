package server

import (
	"github.com/MatticNote/MatticNote/server/account"
	"github.com/MatticNote/MatticNote/server/api"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{})
	})

	app.Get("/@:username", getUserProfile)

	account.ConfigureRoute(app.Group("/account"))
	api.ConfigureRoute(app.Group("/api"))
}
