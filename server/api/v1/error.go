package v1

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func v1BadRequest(c *fiber.Ctx, reason ...string) error {
	if len(reason) > 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":   "BAD_REQUEST",
				"detail": reason[0],
			},
		})
	} else {
		c.Status(http.StatusBadRequest)
		return nil
	}
}

func v1Forbidden(c *fiber.Ctx, reason ...string) error {
	if len(reason) > 0 {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": fiber.Map{
				"code":   "FORBIDDEN",
				"detail": reason[0],
			},
		})
	} else {
		c.Status(http.StatusForbidden)
		return nil
	}
}

func v1NotFound(c *fiber.Ctx, reason ...string) error {
	if len(reason) > 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": fiber.Map{
				"code":   "NOT_FOUND",
				"detail": reason[0],
			},
		})
	} else {
		c.Status(http.StatusNotFound)
		return nil
	}
}

func rateLimitReached(c *fiber.Ctx) error {
	return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
		"error": fiber.Map{
			"code":   "RATE_LIMITED",
			"detail": "The acceptable request limit has been reached",
		},
	})
}
