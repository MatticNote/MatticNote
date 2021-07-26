package api

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/api/graphql"
	"github.com/friendsofgo/graphiql"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func ConfigureRoute(r fiber.Router) {
	//v1.ConfigureRoute(r.Group("/v1"))

	graphqlGroup := r.Group("/graphql")
	graphqlGroup.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "POST",
		AllowHeaders:     "Authorization, Accept, Content-Type, Origin",
		AllowCredentials: false,
		MaxAge:           3600,
	}))
	graphqlGroup.Use(
		internal.RegisterFiberJWT("header", false),
		internal.AuthenticationUser,
		limiter.New(limiter.Config{
			Max:          6000,
			KeyGenerator: internal.RateLimitKeyGen("APIv1"),
			Expiration:   15 * time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusTooManyRequests)
			},
			Storage: config.GetFiberRedisMemory(),
		}),
	)
	graphqlGroup.Post("", graphql.GQLEndpoint)

	graphiQLHandler, err := graphiql.NewGraphiqlHandler("/api/graphql")
	if err != nil {
		panic(err)
	}

	r.Get("/graphiql", adaptor.HTTPHandler(graphiQLHandler))
}
