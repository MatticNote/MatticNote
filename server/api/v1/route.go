package v1

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func ConfigureRoute(r fiber.Router) {
	r.Use(
		internal.RegisterFiberJWT("header", false),
		authenticationUser,
		limiter.New(limiter.Config{
			Max:          6000,
			KeyGenerator: rateLimitKeyGen("APIv1"),
			Expiration:   15 * time.Minute,
			LimitReached: rateLimitReached,
			Storage:      config.GetFiberRedisMemory(),
		}),
	)

	user := r.Group("/user")
	user.Get("/:uuid", getUser)

	note := r.Group("/note")
	note.Post("/", postNote)
	note.Get("/:uuid", getNote)
}
