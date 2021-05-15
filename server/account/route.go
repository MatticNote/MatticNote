package account

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"math/rand"
	"net/http"
)

var tokenCharset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

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

func genToken(size uint8) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = tokenCharset[rand.Intn(len(tokenCharset))]
	}
	return string(b)
}

func ConfigureRoute(r fiber.Router) {
	r.Use(csrf.New(csrf.Config{
		Next:           nil,
		KeyLookup:      fmt.Sprintf("form:%s", csrfFormName),
		CookieName:     "_csrf",
		CookiePath:     "/account",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		ContextKey:     csrfContextKey,
		ErrorHandler:   csrfErrorView,
		KeyGenerator: func() string {
			return genToken(32)
		},
	}))

	r.Get("/register", registerUserView)
	r.Post("/register", registerPost)
}
