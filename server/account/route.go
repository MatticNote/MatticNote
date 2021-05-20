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
		Storage: config.GetFiberRedisMemory(),
	}))

	r.Get("/register", registerUserGet)
	r.Post("/register",
		limiter.New(limiter.Config{
			Next: func(_ *fiber.Ctx) bool {
				return config.Config.Server.DisableAccountRegistrationLimit
			},
			Max: int(config.Config.Server.AccountRegistrationLimitCount),
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_ACCTREG-%s", c.IP())
			},
			Expiration: 24 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		registerPost,
	)

	r.Get("/login", loginUserGet)
	r.Post("/login",
		limiter.New(limiter.Config{
			Max: 30,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_LOGIN-%s", c.IP())
			},
			Expiration: 30 * time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		loginPost,
	)

	r.Get("/logout", destroySession)

	r.Get("/forgot", forgotPasswordGet)
	r.Post("/forgot",
		limiter.New(limiter.Config{
			Max: 10,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_FORGOT-%s", c.IP())
			},
			Expiration: 1 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		forgotPasswordPost,
	)

	r.Get("/issue_confirm_mail", issueConfirmGet)
	r.Post("/issue_confirm_mail",
		limiter.New(limiter.Config{
			Max: 10,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_ISSCM-%s", c.IP())
			},
			Expiration: 1 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				return registerUserView(c, "Rate limit reached")
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		issueConfirmPost,
	)

	r.Get("/verify/:token",
		limiter.New(limiter.Config{
			Max: 30,
			KeyGenerator: func(c *fiber.Ctx) string {
				return fmt.Sprintf("MN_VFTK-%s", c.IP())
			},
			Expiration: 1 * time.Hour,
			LimitReached: func(c *fiber.Ctx) error {
				c.Status(http.StatusTooManyRequests)
				return c.Send([]byte(""))
			},
			Storage: config.GetFiberRedisMemory(),
		}),
		verifyMail,
	)
}
