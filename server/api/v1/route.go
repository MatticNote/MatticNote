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
		internal.AuthenticationUser,
		limiter.New(limiter.Config{
			Max:          6000,
			KeyGenerator: internal.RateLimitKeyGen("APIv1"),
			Expiration:   15 * time.Minute,
			LimitReached: rateLimitReached,
			Storage:      config.GetFiberRedisMemory(),
		}),
	)

	user := r.Group("/user")
	user.Get("/:uuid", getUser)
	user.Post("/:uuid/follow", followUser)
	user.Delete("/:uuid/follow", unfollowUser)
	user.Post("/:uuid/block", blockUser)
	user.Delete("/:uuid/block", unblockUser)

	userFollowRequests := r.Group("/follow_request")
	userFollowRequests.Get("/", listRequests)
	userFollowRequests.Post("/:uuid/accept", acceptRequests)
	userFollowRequests.Post("/:uuid/reject", rejectRequests)

	note := r.Group("/note")
	note.Post("/", postNote)
	note.Get("/:uuid", getNote)
}
