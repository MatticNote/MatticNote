package oauth

import (
	"github.com/MatticNote/MatticNote/internal/oauth"
	"net/http"
)

func revoke(w http.ResponseWriter, r *http.Request) {
	err := oauth.Server.NewRevocationRequest(r.Context(), r)
	oauth.Server.WriteRevocationResponse(w, err)
}
