package account

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"net/http"
	"time"
)

const (
	csrfContextKey = "CSRF"
	csrfFormName   = "_csrf"
)

func csrfErrorView(c *fiber.Ctx, _ error) error {
	return c.Status(http.StatusForbidden).Render(
		"csrf_failed",
		fiber.Map{},
	)
}

func ConfigureRoute(r fiber.Router) {
	r.Use(csrf.New(csrf.Config{
		Next:           nil,
		KeyLookup:      fmt.Sprintf("form:%s", csrfFormName),
		CookieName:     "_csrf",
		CookiePath:     "/account",
		CookieSecure:   config.Config.Server.CookieSecure,
		CookieHTTPOnly: true,
		ContextKey:     csrfContextKey,
		ErrorHandler:   csrfErrorView,
		KeyGenerator: func() string {
			return misc.GenToken(32)
		},
	}))

	r.Get("/register", registerUserGet)
	r.Post("/register",
		limiter.New(limiter.Config{
			Next: func(_ *fiber.Ctx) bool {
				return config.Config.Server.DisableAccountRegistrationLimit
			},
			Max: int(config.Config.Server.AccountRegistrationLimitCount),
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			Expiration: 24 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
		}),
		registerPost,
	)

	r.Get("/login", loginUserGet)
	r.Post("/login",
		limiter.New(limiter.Config{
			Max: 30,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			Expiration: 30 * time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
		}),
		loginPost,
	)

	r.Get("/logout", destroySession)

	r.Get("/forgot", forgotPasswordGet)
	r.Post("/forgot",
		limiter.New(limiter.Config{
			Max: 10,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			Expiration: 1 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
		}),
		forgotPasswordPost,
	)

}
