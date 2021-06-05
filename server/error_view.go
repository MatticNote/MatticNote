package server

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func NotFoundView(c *fiber.Ctx) error {
	if c.Accepts("html") != "" {
		return c.Status(fiber.StatusNotFound).Render(
			"404",
			fiber.Map{},
		)
	} else {
		c.Status(fiber.StatusNotFound)
		return nil
	}
}

func GoneView(c *fiber.Ctx) error {
	c.Status(fiber.StatusGone)
	return nil
}

func ErrorView(c *fiber.Ctx, err error) error {
	switch err {
	case fiber.ErrUnauthorized:
		c.Status(fiber.StatusUnauthorized)
	case fiber.ErrForbidden:
		c.Status(fiber.StatusForbidden)
	case fiber.ErrNotFound:
		return NotFoundView(c)
	case fiber.ErrGone:
		return GoneView(c)
	default:
		return c.Status(http.StatusInternalServerError).Render(
			"5xx",
			fiber.Map{
				"Error": err.Error(),
			},
		)
	}

	return nil
}
