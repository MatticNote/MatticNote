package account

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func registerUserView(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).Render(
		"register",
		fiber.Map{
			"CSRFFormName": csrfFormName,
			"CSRFToken":    c.Context().UserValue(csrfContextKey).(string),
		},
		"layout/account",
	)
}

func registerPost(c *fiber.Ctx) error {
	return c.Status(200).SendString("POST")
}
