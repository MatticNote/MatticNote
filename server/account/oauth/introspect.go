package oauth

import (
	"github.com/MatticNote/MatticNote/internal/oauth"
	"github.com/ory/fosite"
	"net/http"
)

func introspect(w http.ResponseWriter, r *http.Request) {
	session := &fosite.DefaultSession{}
	res, err := oauth.Server.NewIntrospectionRequest(r.Context(), r, session)
	if err != nil {
		oauth.Server.WriteIntrospectionError(w, err)
		return
	}
	oauth.Server.WriteIntrospectionResponse(w, res)
}
