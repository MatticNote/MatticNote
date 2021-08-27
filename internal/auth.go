package internal

import (
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/MatticNote/MatticNote/internal/user"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strings"
)

const (
	LoginUserLocal = "loginUser"
)

func AuthenticationUser(c *fiber.Ctx) error {
	// JWT token authentication
	headerSplit := strings.Split(c.Get(signature.AuthHeaderName, ""), " ")
	if len(headerSplit) > 0 && strings.TrimSpace(headerSplit[0]) == signature.AuthSchemeName {
		token, ok := c.Locals(signature.JWTContextKey).(*jwt.Token)
		if ok {
			claim := token.Claims.(jwt.MapClaims)
			usr, err := user.GetLocalUser(uuid.MustParse(claim["sub"].(string)))
			if err == nil {
				c.Locals(LoginUserLocal, usr)
			}
		} else {
			return fiber.ErrUnauthorized
		}
	}

	return c.Next()
}
