package setting

import (
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/gofiber/fiber/v2"
)

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
