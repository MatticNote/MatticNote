package v1

import "github.com/gofiber/fiber/v2"

type (
	apiErrorJson struct {
		Code   string `json:"code"`
		Detail string `json:"detail"`
	}
)

func apiError(c *fiber.Ctx, status int, code, detail string) error {
	return c.Status(status).JSON(apiErrorJson{
		Code:   code,
		Detail: detail,
	})
}

func apiBadRequest(c *fiber.Ctx, detail string) error {
	return apiError(c, fiber.StatusBadRequest, "BAD_REQUEST", detail)
}

func apiUnauthorized(c *fiber.Ctx) error {
	return apiError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Invalid credentials.")
}

func apiNotFound(c *fiber.Ctx, detail string) error {
	return apiError(c, fiber.StatusNotFound, "NOT_FOUND", detail)
}

func apiGone(c *fiber.Ctx, detail string) error {
	return apiError(c, fiber.StatusGone, "GONE", detail)
}
