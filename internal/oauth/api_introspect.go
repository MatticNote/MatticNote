package oauth

import (
	"context"
	"github.com/MatticNote/MatticNote/internal/ist"
	"github.com/MatticNote/MatticNote/internal/user"
	"github.com/google/uuid"
	"github.com/ory/fosite"
)

func APIIntrospect(token string) (*ist.LocalUserStruct, string, error) {
	_, ar, err := Server.IntrospectToken(context.Background(), token, fosite.AccessToken, &fosite.DefaultSession{})
	if err != nil {
		return nil, "", err
	}
	targetUuid, err := uuid.Parse(ar.GetSession().GetUsername())
	if err != nil {
		return nil, "", err
	}
	localUser, err := user.GetLocalUser(targetUuid)
	if err != nil {
		return nil, "", err
	}
	return localUser, ar.GetClient().GetID(), err
}
