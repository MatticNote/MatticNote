package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func ErrorView(c *fiber.Ctx, err error) error {
	switch err {
	case fiber.ErrNotFound:
		return NotFoundView(c)
	case fiber.ErrForbidden:
		return ForbiddenView(c)
	case fiber.ErrUnprocessableEntity:
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	default:
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("%T: %s", err, err.Error()))
	}
}

func NotFoundView(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).Render(
		"error/404",
		fiber.Map{},
		"error/_layout",
	)
}

func ForbiddenView(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).Render(
		"error/403",
		fiber.Map{},
		"error/_layout",
	)
}
