package server

import (
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/server/account"
	"github.com/MatticNote/MatticNote/server/api"
	"github.com/MatticNote/MatticNote/server/common"
	"github.com/MatticNote/MatticNote/server/setting"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/segmentio/ksuid"
	"strings"
	"time"
)

func ConfigureRoute(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{})
	})

	app.Get("/@:username", getUserProfile)
	app.Get("/@:username/:noteId", getUserNote)

	app.Use(csrf.New(csrf.Config{
		KeyLookup: fmt.Sprintf("form:%s", common.CSRFFormName),
		KeyGenerator: func() string {
			return ksuid.New().String()
		},
		Next: func(c *fiber.Ctx) bool {
			var path = c.Path()
			switch {
			case strings.HasPrefix(path, "/account"):
				return false
			case strings.HasPrefix(path, "/settings"):
				return false
			}
			return true
		},
		Expiration:        15 * time.Minute,
		CookiePath:        "/",
		CookieName:        common.CSRFCookieName,
		ContextKey:        common.CSRFContextKey,
		CookieHTTPOnly:    true,
		CookieSessionOnly: true,
		CookieSameSite:    fiber.CookieSameSiteStrictMode,
		Storage:           database.FiberStorage,
		ErrorHandler:      common.CSRFErrorHandler,
	}))

	account.ConfigureRoute(app.Group("/account"))
	api.ConfigureRoute(app.Group("/api"))
	setting.ConfigureRoute(app.Group("/settings", common.ValidateCookie))
}
