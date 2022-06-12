package server

import "github.com/gofiber/fiber/v2"

func ErrorView(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
}
