package oauth

import (
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func ConfigureRoute(r fiber.Router) {
	r.All("/authorize",
		signature.RegisterFiberJWT("cookie", true),
		authorize,
	)
	r.All("/token",
		authorizeToken,
	)
	r.All("/revoke",
		adaptor.HTTPHandlerFunc(revoke),
	)
}
