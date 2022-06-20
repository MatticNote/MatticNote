package account

import (
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	ia "github.com/MatticNote/MatticNote/internal/account"
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

func validateCookie(c *fiber.Ctx) error {
	token := c.Cookies(ia.TokenCookieName)
	if token == "" {
		return c.Redirect("/account/login")
	}

	user, err := ia.GetUserFromToken(token)
	if err != nil {
		switch err {
		case ia.ErrUserNotFound, ia.ErrInvalidUserToken, ia.ErrUserGone:
			return c.Redirect("/account/logout")
		default:
			return err
		}
	}

	c.Locals("currentUser", user)

	return c.Next()
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

	r.Get("/login", loginGet)
	r.Post("/login", loginPost)

	r.Get("/logout", logout)

	r.Get("/register", registerGet)
	r.Post("/register", registerPost)

	r.Get("/verify/:token", verifyEmailToken)

	r.Get("/register-username", validateCookie, registerUsernameGet)
	r.Post("/register-username", validateCookie, registerUsernamePost)
}
