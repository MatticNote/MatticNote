package oauth

import (
	"github.com/google/uuid"
	"github.com/ory/fosite"
	"golang.org/x/crypto/bcrypt"
)

type mnOAuthClient struct {
	ID           string
	Secret       []byte
	Name         string
	Owner        *uuid.UUID
	RedirectURIs []string
	Scopes       []string
}

func (c *mnOAuthClient) GetID() string {
	return c.ID
}

func (c *mnOAuthClient) GetHashedSecret() []byte {
	hashed, err := bcrypt.GenerateFromPassword(c.Secret, 10)
	if err != nil {
		panic(err)
	}
	return hashed
}

func (c *mnOAuthClient) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *mnOAuthClient) GetGrantTypes() fosite.Arguments {
	return fosite.Arguments{"authorization_code", "refresh_token"}
}

func (c *mnOAuthClient) GetResponseTypes() fosite.Arguments {
	return fosite.Arguments{"code"}
}

func (c *mnOAuthClient) GetScopes() fosite.Arguments {
	return c.Scopes
}

func (c *mnOAuthClient) IsPublic() bool {
	return true
}

func (c *mnOAuthClient) GetAudience() fosite.Arguments {
	return nil
}
