package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func ConfigureRoute(r fiber.Router) {
	r.Use(limiter.New(limiter.Config{
		Max: 6000,
		KeyGenerator: func(c *fiber.Ctx) string {
			// TODO: ヘッダーのユーザトークンを基にレートリミットできるようにする
			return c.IP()
		},
		Expiration:   15 * time.Minute,
		LimitReached: v1TooManyRequests,
	}))

	user := r.Group("/user")

	user.Get(":uuid", getUser)
}
