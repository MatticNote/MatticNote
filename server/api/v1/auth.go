package v1

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"strings"
)

const (
	loginUserLocal = "loginUser"
)

func authenticationUser(c *fiber.Ctx) error {

	// JWT token authentication
	headerSplit := strings.Split(c.Get(internal.AuthHeaderName, ""), " ")
	if len(headerSplit) > 0 && strings.TrimSpace(headerSplit[0]) == internal.AuthSchemeName {
		token, ok := c.Locals(internal.JWTContextKey).(*jwt.Token)
		if ok {
			claim := token.Claims.(jwt.MapClaims)
			usr, err := internal.GetLocalUser(claim["sub"].(string))
			if err == nil {
				c.Locals(loginUserLocal, usr)
			}
		} else {
			return fiber.ErrUnauthorized
		}
	}

	return c.Next()
}
