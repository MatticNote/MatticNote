package account

import (
	"github.com/MatticNote/MatticNote/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/segmentio/ksuid"
	"time"
)

const csrfContextKey = "csrf_token"

func ConfigureRoute(r fiber.Router) {
	r.Use(csrf.New(csrf.Config{
		KeyLookup: "form:csrf_token",
		KeyGenerator: func() string {
			return ksuid.New().String()
		},
		Expiration:        15 * time.Minute,
		ContextKey:        csrfContextKey,
		CookieHTTPOnly:    true,
		CookieSessionOnly: true,
		CookieSameSite:    fiber.CookieSameSiteStrictMode,
		Storage:           database.FiberStorage,
	}))

	r.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("account/login", fiber.Map{
			"csrf_token": c.Locals(csrfContextKey),
		})
	})

	r.Get("/register", func(c *fiber.Ctx) error {
		return c.Render("account/register", fiber.Map{
			"csrf_token": c.Locals(csrfContextKey),
		})
	})
}
