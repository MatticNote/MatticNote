package oauth

import (
	"github.com/MatticNote/MatticNote/internal/oauth"
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/MatticNote/MatticNote/mn_template"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/ory/fosite"
	"log"
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
			tmpl, err := mn_template.LoadOAuthTemplate()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = tmpl.ExecuteTemplate(w, "authorize.html", fiber.Map{
				"name":            ar.GetClient().GetID(),
				"scopes":          ar.GetRequestedScopes(),
				"csrf_form_name":  misc.CSRFFormName,
				"csrf_form_token": csrfToken,
			})
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		if r.Form.Get("authorize") != "grant" {
			oauth.Server.WriteAuthorizeError(w, ar, &fosite.RFC6749Error{
				DescriptionField: "Declined by user.",
				ErrorField:       "access_denied",
				CodeField:        http.StatusForbidden,
			})
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
