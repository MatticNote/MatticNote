package account

import (
	"github.com/MatticNote/MatticNote/server/common"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(r fiber.Router) {
	common.InitAccountSessionStore()

	r.Get("/login", loginGet)
	r.Post("/login", loginPost)

	r.Get("/logout", logout)

	r.Get("/register", registerGet)
	r.Post("/register", registerPost)

	r.Get("/reset-password", resetPasswordNewGet)
	r.Post("/reset-password", resetPasswordNewPost)
	r.Get("/reset-password/:token", resetPasswordFormGet)
	r.Post("/reset-password/:token", resetPasswordFormPost)

	r.Get("/verify/:token", verifyEmailToken)

	r.Get("/register-username", common.ValidateCookie, registerUsernameGet)
	r.Post("/register-username", common.ValidateCookie, registerUsernamePost)
}
