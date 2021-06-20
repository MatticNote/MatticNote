package settings

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(r fiber.Router) {
	r.Use(internal.RegisterFiberJWT("cookie", true))

	r.Get("/profile", editProfileGet)
	r.Post("/profile", editProfilePost)

	r.Get("/password", editPasswordGet)
	r.Post("/password", editPasswordPost)

	r.Get("/security", securityPageGet)
	r.Get("/security/2fa/setup", setup2faGet)
	r.Post("/security/2fa/setup", setup2faPost)
	r.Get("/security/2fa/backup", get2faBackup)
	r.Post("/security/2fa/backup", regenerate2faBackup)
}
