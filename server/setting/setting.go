package setting

import (
	"github.com/gofiber/fiber/v2"
)

type (
	actionFormStruct struct {
		Action string
	}
)

func ConfigureRoute(r fiber.Router) {
	r.Get("/core", settingCoreGet)
	r.Post("/core", settingCorePost)

	r.Get("/security", settingSecurityGet)

	r.Get("/invite", settingInviteGet)
	r.Post("/invite", settingInvitePost)
}
