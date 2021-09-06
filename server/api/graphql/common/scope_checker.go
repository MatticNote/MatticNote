package common

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal/auth"
	"github.com/MatticNote/MatticNote/internal/oauth"
	"github.com/graphql-go/graphql"
	"github.com/ory/fosite"
)

func ScopeCheck(p graphql.ResolveParams, requiredScopes ...string) error {
	if method, ok := p.Context.Value(ContextAuthorizeMethod).(auth.AuthorizeMethod); ok && method == auth.JWT {
		return nil
	}

	if token, ok := p.Context.Value(ContextOAuthToken).(string); ok {
		err := oauth.ScopeIntrospect(token, requiredScopes...)
		if err != nil {
			switch {
			case errors.Is(err, fosite.ErrInvalidScope):
				return ErrNotEnoughScopes
			default:
				return err
			}
		}
		return nil
	} else {
		return ErrNotEnoughScopes
	}
}
