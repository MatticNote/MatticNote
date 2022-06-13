package server

import (
	"github.com/MatticNote/MatticNote/server/account"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("hello", fiber.Map{})
	})
	account.ConfigureRoute(app.Group("/account"))
}