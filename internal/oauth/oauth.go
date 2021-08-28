package oauth

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
	"github.com/ory/fosite/token/jwt"
	"time"
)

var fositeCfg = &compose.Config{
	AccessTokenLifespan:  3 * time.Hour,
	RefreshTokenLifespan: -1,
	HashCost:             12,
	RefreshTokenScopes:   []string{},
}

var Server fosite.OAuth2Provider

func InitOAuth() {
	store = storage.NewExampleStore()
	Server = compose.Compose(
		fositeCfg,
		store,
		&compose.CommonStrategy{
			CoreStrategy: compose.NewOAuth2HMACStrategy(
				fositeCfg,
				[]byte(config.Config.Server.OAuthSecretKey),
				nil,
			),
			JWTStrategy: &jwt.RS256JWTStrategy{
				PrivateKey: signature.GetPrivateKey(),
			},
		},
		nil,

		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		//compose.RFC7523AssertionGrantFactory,

		compose.OAuth2TokenIntrospectionFactory,
		compose.OAuth2TokenRevocationFactory,
		//compose.OAuth2PKCEFactory,
	)
}
