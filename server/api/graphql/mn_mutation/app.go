package mn_mutation

import (
	"github.com/MatticNote/MatticNote/internal/oauth"
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_type"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var CreateApp = &graphql.Field{
	Name:        "CreateApp",
	Description: "Create a app.",
	Args: graphql.FieldConfigArgument{
		"name": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Application name. required.",
		},
		"scope": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.NewList(graphql.String)),
			Description: "Scope. required.",
		},
		"redirectUris": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.NewList(graphql.String)),
			Description: "Redirect uris. required.",
		},
	},
	Type: graphql.NewNonNull(mn_type.NewAppQLType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var ownerUuid *uuid.UUID
		if currentUser, ok := p.Context.Value("currentUser").(uuid.UUID); ok {
			ownerUuid = &currentUser
		}

		redirectUrisRaw := p.Args["redirectUris"].([]interface{})
		var redirectUris []string
		for _, val := range redirectUrisRaw {
			redirectUris = append(redirectUris, val.(string))
		}
		scopesRaw := p.Args["scope"].([]interface{})
		var scopes []string
		for _, val := range scopesRaw {
			scopes = append(scopes, val.(string))
		}

		newCliI, err := oauth.NewClient(p.Args["name"].(string), ownerUuid, redirectUris, scopes)
		if err != nil {
			return nil, err
		}
		newCli := newCliI.(*oauth.MNOAuthClient)

		return mn_type.NewAppType{
			ClientID:     newCli.GetID(),
			SecretID:     string(newCli.Secret),
			Name:         newCli.Name,
			Scopes:       newCli.GetScopes(),
			RedirectUris: newCli.GetRedirectURIs(),
			Owner:        newCli.Owner,
		}, nil
	},
}
