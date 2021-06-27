package well_known

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ConfigureRoute(r fiber.Router) {
	r.Get("/change-password", redirectChPasswd)
	r.Get("/webfinger",
		cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowMethods:     "GET",
			AllowHeaders:     "Accept, Content-Type, Origin",
			AllowCredentials: false,
			MaxAge:           3600,
		}),
		webfinger,
	)
}
