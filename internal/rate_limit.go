package internal

import (
	"fmt"
	"github.com/MatticNote/MatticNote/internal/auth"
	"github.com/MatticNote/MatticNote/internal/ist"
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func RateLimitKeyGen(prefix string) func(c *fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		headerSplit := strings.Split(c.Get(signature.AuthHeaderName, ""), " ")
		if len(headerSplit) > 0 && strings.TrimSpace(headerSplit[0]) == signature.AuthSchemeName {
			currentUsr := c.Locals(auth.LoginUserLocal).(*ist.LocalUserStruct)
			return fmt.Sprintf("MN_%s-jwt@%s", prefix, currentUsr.Uuid)
		} else if val, ok := c.Locals(auth.AuthorizeMethodLocal).(auth.AuthorizeMethod); ok && val == auth.OAuth {
			currentUsr := c.Locals(auth.LoginUserLocal).(*ist.LocalUserStruct)
			clientId := c.Locals(auth.OAuthClientIDLocal).(string)
			return fmt.Sprintf("MN_%s-%s@%s", prefix, clientId, currentUsr.Uuid)
		} else {
			return fmt.Sprintf("MN_%s-%s", prefix, c.IP())
		}
	}
}
