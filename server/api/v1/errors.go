package v1

import "github.com/gofiber/fiber/v2"

type (
	apiErrorJson struct {
		Code   string `json:"code"`
		Detail string `json:"detail"`
	}
)

func apiBadRequest(c *fiber.Ctx, detail string) error {
	return c.Status(fiber.StatusBadRequest).JSON(apiErrorJson{
		Code:   "BAD_REQUEST",
		Detail: detail,
	})
}

func apiUnauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(apiErrorJson{
		Code:   "UNAUTHORIZED",
		Detail: "Invalid credentials.",
	})
}

func apiNotFound(c *fiber.Ctx, detail string) error {
	return c.Status(fiber.StatusNotFound).JSON(apiErrorJson{
		Code:   "NOT_FOUND",
		Detail: detail,
	})
}

func apiGone(c *fiber.Ctx, detail string) error {
	return c.Status(fiber.StatusGone).JSON(apiErrorJson{
		Code:   "GONE",
		Detail: detail,
	})
}
