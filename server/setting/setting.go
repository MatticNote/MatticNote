package setting

import (
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(r fiber.Router) {
	r.Get("/core", settingCoreGet)
	r.Get("/security", settingSecurityGet)
}

func settingCoreGet(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)

	isEmailVerified, err := account.IsEmailVerified(user.ID)
	if err != nil {
		return err
	}
	email, err := account.GetUserEmail(user.ID)
	if err != nil {
		return err
	}

	return c.Render(
		"setting/core",
		fiber.Map{
			"user":                user,
			"isUserEmailVerified": isEmailVerified,
			"userEmail":           email,
		},
		"setting/_layout",
	)
}

func settingSecurityGet(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)
	token := c.Cookies(account.TokenCookieName)

	tokenList, err := account.ListUserToken(user.ID)
	if err != nil {
		return err
	}

	return c.Render(
		"setting/security",
		fiber.Map{
			"user":         user,
			"tokenList":    tokenList,
			"currentToken": token,
		},
		"setting/_layout",
	)
}
