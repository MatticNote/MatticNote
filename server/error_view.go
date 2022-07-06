package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func ErrorView(c *fiber.Ctx, err error) error {
	switch err {
	case fiber.ErrNotFound:
		return NotFoundView(c)
	default:
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("%T: %s", err, err.Error()))
	}
}

func NotFoundView(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).Render(
		"404",
		fiber.Map{},
	)
}
