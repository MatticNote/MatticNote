package server

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/account"
	"github.com/MatticNote/MatticNote/server/api"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

var emptyBody = []byte("")

func ConfigureRoute(app *fiber.App) {
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

	account.ConfigureRoute(app.Group("/account"))
	api.ConfigureRoute(app.Group("/api"))
}

// internal views

func NotFoundView(c *fiber.Ctx) error {
	if c.Accepts("html") != "" {
		return c.Status(http.StatusNotFound).Render(
			"404",
			fiber.Map{},
		)
	} else {
		return c.Status(http.StatusNotFound).Send(emptyBody)
	}
}

func ErrorView(c *fiber.Ctx, err error) error {
	switch err {
	case fiber.ErrUnauthorized:
		c.Status(fiber.StatusUnauthorized)
	case fiber.ErrForbidden:
		c.Status(fiber.StatusForbidden)
	default:
		return c.Status(http.StatusInternalServerError).Render(
			"5xx",
			fiber.Map{
				"Error": err.Error(),
			},
		)
	}

	return c.Send(emptyBody) // empty body
}
