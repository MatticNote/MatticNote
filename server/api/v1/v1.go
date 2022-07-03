package v1

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func ConfigureRoute(r fiber.Router) {
	r.Use(func(c *fiber.Ctx) error {
		authHeader := strings.SplitN(c.Get("Authorization"), " ", 2)
		if len(authHeader) == 2 {
			switch strings.ToLower(authHeader[0]) {
			case "bearer":
			case "token":
				if authHeader[1] != "" {
					user, err := account.GetUserFromToken(authHeader[1])
					if err != nil {
						if errors.Is(err, account.ErrInvalidUserToken) {
							return apiUnauthorized(c)
						} else {
							return err
						}
					}
					if user != nil {
						c.Locals("currentUser", user)
					}
				}
			}
		}
		return c.Next()
	})

	userApiRoute(r.Group("/users"))
	noteApiRoute(r.Group("/notes"))
}
