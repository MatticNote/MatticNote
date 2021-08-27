package oauth

import (
	"fmt"
	"github.com/MatticNote/MatticNote/internal/oauth"
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/ory/fosite"
	"net/http"
)

func authorize(c *fiber.Ctx) error {
	jwtCurrentUserKey := c.Locals(signature.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrUnauthorized
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)
	sub, ok := claim["sub"].(string)
	if !ok {
		return fiber.ErrUnauthorized
	}

	csrfToken := c.Locals(misc.CSRFContextKey).(string)

	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ar, err := oauth.Server.NewAuthorizeRequest(r.Context(), r)
		if err != nil {
			oauth.Server.WriteAuthorizeError(w, ar, err)
			return
		}

		mySessionData := &fosite.DefaultSession{
			Username: sub,
		}

		if r.Method != http.MethodPost {
			var requestedScopes string
			for _, this := range ar.GetRequestedScopes() {
				requestedScopes += fmt.Sprintf(`<li>%s</li>`, this)
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(`<h1>Login page</h1>`))
			_, _ = w.Write([]byte(fmt.Sprintf(`
			<p>Howdy! This is the log in page. For this example, it is enough to supply the username.</p>
			<form method="post">
				By logging in, you consent to grant these scopes:
				<ul>%s</ul>
				<input type="hidden" name="%s" value="%s">
				<input type="submit">
			</form>
		`, requestedScopes,
				misc.CSRFFormName,
				csrfToken,
			)))
			return
		}

		for _, scope := range ar.GetRequestedScopes() {
			ar.GrantScope(scope)
		}

		res, err := oauth.Server.NewAuthorizeResponse(c.Context(), ar, mySessionData)
		if err != nil {
			oauth.Server.WriteAuthorizeError(w, ar, err)
			return
		}

		oauth.Server.WriteAuthorizeResponse(w, ar, res)
	})(c)
}
