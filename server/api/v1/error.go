package v1

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func renderError(c *fiber.Ctx, status int, code string, reason ...string) error {
	if len(reason) > 0 {
		return c.Status(status).JSON(fiber.Map{
			"error": fiber.Map{
				"code":   code,
				"detail": reason[0],
			},
		})
	} else {
		c.Status(status)
		return nil
	}
}

func badRequest(c *fiber.Ctx, reason ...string) error {
	return renderError(c, http.StatusBadRequest, "BAD_REQUEST", reason...)
}

func unauthorized(c *fiber.Ctx, reason ...string) error {
	return renderError(c, http.StatusUnauthorized, "UNAUTHORIZED", reason...)
}

func forbidden(c *fiber.Ctx, reason ...string) error {
	return renderError(c, http.StatusForbidden, "FORBIDDEN", reason...)
}

func notFound(c *fiber.Ctx, reason ...string) error {
	return renderError(c, http.StatusNotFound, "NOT_FOUND", reason...)
}

func rateLimitReached(c *fiber.Ctx) error {
	return renderError(c, http.StatusTooManyRequests, "RATE_REACHED", "The acceptable request limit has been reached")
}
