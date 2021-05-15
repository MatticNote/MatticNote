package api

import (
	v1 "github.com/MatticNote/MatticNote/server/api/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ConfigureRoute(r fiber.Router) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Authorization, Accept, Content-Type, Origin",
		AllowCredentials: false,
		MaxAge:           3600,
	}))

	v1.ConfigureRoute(r.Group("/v1"))
}
