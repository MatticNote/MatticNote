package v1

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func rateLimitKeyGen(c *fiber.Ctx) string {
	// TODO: ヘッダーのユーザトークンを基にレートリミットできるようにする
	return fmt.Sprintf("MN_APIv1-%s", c.IP())
}

func ConfigureRoute(r fiber.Router) {
	r.Use(limiter.New(limiter.Config{
		Max:          6000,
		KeyGenerator: rateLimitKeyGen,
		Expiration:   15 * time.Minute,
		LimitReached: v1TooManyRequests,
		Storage:      config.GetFiberRedisMemory(),
	}))

	user := r.Group("/user")

	user.Get(":uuid", getUser)
}
