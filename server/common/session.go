package common

import (
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

const (
	AccountSessionCookieName = "mn_as"
	AccountSessionRedirectTo = "redirectTo"
)

var AccountSession *session.Store

func InitAccountSessionStore() {
	AccountSession = session.New(session.Config{
		Expiration:     15 * time.Minute,
		Storage:        database.FiberStorage,
		CookiePath:     "/account",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		CookieSameSite: fiber.CookieSameSiteStrictMode,
		KeyGenerator: func() string {
			return internal.GenerateRandString(128)
		},
		KeyLookup: fmt.Sprintf("cookie:%s", AccountSessionCookieName),
	})
}
