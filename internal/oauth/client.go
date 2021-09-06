package oauth

import (
	"context"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/google/uuid"
	"github.com/ory/fosite"
	"golang.org/x/crypto/bcrypt"
)

type MNOAuthClient struct {
	ID           string
	Secret       []byte
	Name         string
	Owner        *uuid.UUID
	RedirectURIs []string
	Scopes       []string
}

func (c *MNOAuthClient) GetID() string {
	return c.ID
}

func (c *MNOAuthClient) GetHashedSecret() []byte {
	hashed, err := bcrypt.GenerateFromPassword(c.Secret, 12)
	if err != nil {
		panic(err)
	}
	return hashed
}

func (c *MNOAuthClient) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *MNOAuthClient) GetGrantTypes() fosite.Arguments {
	return fosite.Arguments{"authorization_code", "refresh_token"}
}

func (c *MNOAuthClient) GetResponseTypes() fosite.Arguments {
	return fosite.Arguments{"code"}
}

func (c *MNOAuthClient) GetScopes() fosite.Arguments {
	return c.Scopes
}

func (c *MNOAuthClient) IsPublic() bool {
	return true
}

func (c *MNOAuthClient) GetAudience() fosite.Arguments {
	return nil
}

func NewClient(name string, owner *uuid.UUID, redirectURIs, scopes []string) (fosite.Client, error) {
	newID := misc.GenToken(32)
	newSecret := []byte(misc.GenToken(64))
	_, err := database.DBPool.Exec(
		context.Background(),
		"insert into oauth_client(client_key, client_secret, name, client_owner, redirect_uris, scopes) values ($1, $2, $3, $4, $5, $6);",
		newID,
		newSecret,
		name,
		owner,
		redirectURIs,
		scopes,
	)
	if err != nil {
		return nil, err
	}

	newClient := &MNOAuthClient{
		ID:           newID,
		Secret:       newSecret,
		Name:         name,
		Owner:        owner,
		RedirectURIs: redirectURIs,
		Scopes:       scopes,
	}

	return newClient, nil
}
