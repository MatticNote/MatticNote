package account

import (
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/segmentio/ksuid"
	"time"
)

const (
	csrfFormName   = "csrf_token"
	csrfContextKey = "csrf_token"
)

func csrfErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusForbidden).SendString(fmt.Sprintf("CSRF Error: %s", err.Error()))
}

func ConfigureRoute(r fiber.Router) {
	r.Use(csrf.New(csrf.Config{
		KeyLookup: fmt.Sprintf("form:%s", csrfFormName),
		KeyGenerator: func() string {
			return ksuid.New().String()
		},
		Expiration:        15 * time.Minute,
		ContextKey:        csrfContextKey,
		CookieHTTPOnly:    true,
		CookieSessionOnly: true,
		CookieSameSite:    fiber.CookieSameSiteStrictMode,
		Storage:           database.FiberStorage,
		ErrorHandler:      csrfErrorHandler,
	}))

	r.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("account/login", fiber.Map{
			"csrf_name":  csrfFormName,
			"csrf_token": c.Locals(csrfContextKey),
		})
	})

	r.Get("/register", registerGet)
	r.Post("/register", registerPost)
}
