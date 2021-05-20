package v1

import (
	"fmt"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func rateLimitKeyGen(prefix string) func(c *fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		headerSplit := strings.Split(c.Get(internal.AuthHeaderName, ""), " ")
		if len(headerSplit) > 0 && strings.TrimSpace(headerSplit[0]) == internal.AuthSchemeName {
			currentUsr := c.Locals(loginUserLocal).(*internal.LocalUserStruct)
			return fmt.Sprintf("MN_%s-jwt@%s", prefix, currentUsr.Uuid)
		} else {
			return fmt.Sprintf("MN_%s-%s", prefix, c.IP())
		}
	}
}
