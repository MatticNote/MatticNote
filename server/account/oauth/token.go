package oauth

import (
	"github.com/MatticNote/MatticNote/internal/oauth"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/ory/fosite"
	"net/http"
)

func authorizeToken(c *fiber.Ctx) error {
	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mySessionData := new(fosite.DefaultSession)
		ar, err := oauth.Server.NewAccessRequest(r.Context(), r, mySessionData)
		if err != nil {
			oauth.Server.WriteAccessError(w, ar, err)
			return
		}

		res, err := oauth.Server.NewAccessResponse(r.Context(), ar)
		if err != nil {
			oauth.Server.WriteAccessError(w, ar, err)
		}

		oauth.Server.WriteAccessResponse(w, ar, res)
	})(c)
}
