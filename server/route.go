package server

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/account"
	"github.com/MatticNote/MatticNote/server/ap"
	"github.com/MatticNote/MatticNote/server/api"
	"github.com/MatticNote/MatticNote/server/nodeinfo"
	"github.com/MatticNote/MatticNote/server/well_known"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		// Security Header
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		return c.Next()
	})
	app.Get("/", internal.RegisterFiberJWT("cookie", false), func(c *fiber.Ctx) error {
		_, isLogin := c.Locals(internal.JWTContextKey).(*jwt.Token)

		field := fiber.Map{
			"isLogin": isLogin,
		}

		return c.Render(
			"index",
			field,
			"_layout/index",
		)
	})
	app.Get("/login", func(c *fiber.Ctx) error {
		// alias path
		return c.Redirect("/account/login", fiber.StatusPermanentRedirect)
	})

	app.Get("/@:username", userProfileController)
	app.Get("/user/:username", userProfileController) // alias path

	account.ConfigureRoute(app.Group("/account"))
	api.ConfigureRoute(app.Group("/api"))
	well_known.ConfigureRoute(app.Group("/.well-known"))
	nodeinfo.ConfigureRoute(app.Group("/nodeinfo"))
	ap.ConfigureRoute(app.Group("/activity"))
}
