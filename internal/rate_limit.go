package internal

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func RateLimitKeyGen(prefix string) func(c *fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		headerSplit := strings.Split(c.Get(AuthHeaderName, ""), " ")
		if len(headerSplit) > 0 && strings.TrimSpace(headerSplit[0]) == AuthSchemeName {
			currentUsr := c.Locals(LoginUserLocal).(*LocalUserStruct)
			return fmt.Sprintf("MN_%s-jwt@%s", prefix, currentUsr.Uuid)
		} else {
			return fmt.Sprintf("MN_%s-%s", prefix, c.IP())
		}
	}
}
