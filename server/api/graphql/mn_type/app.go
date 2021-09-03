package mn_type

import (
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

type NewAppType struct {
	ClientID     string
	SecretID     string
	Name         string
	Scopes       []string
	RedirectUris []string
	Owner        *uuid.UUID
}

var NewAppQLType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "NewApp",
	Description: "New application type",
	Fields: graphql.Fields{
		"clientID": &graphql.Field{
			Name:        "ClientID",
			Description: "client id.",
			Type:        graphql.NewNonNull(graphql.ID),
		},
		"secretID": &graphql.Field{
			Name:        "SecretID",
			Description: "Secret id. Do not share with third parties!",
			Type:        graphql.NewNonNull(graphql.String),
		},
		"name": &graphql.Field{
			Name:        "Name",
			Description: "Application name.",
			Type:        graphql.NewNonNull(graphql.String),
		},
		"scopes": &graphql.Field{
			Name:        "Scopes",
			Description: "List of scopes that can be used by this app",
			Type:        graphql.NewNonNull(graphql.NewList(graphql.String)),
		},
		"redirectUris": &graphql.Field{
			Name:        "RedirectUris",
			Description: "Callback destination for this app",
			Type:        graphql.NewNonNull(graphql.NewList(graphql.String)),
		},
		"owner": &graphql.Field{
			Name:        "Owner",
			Description: "Owner's user uuid.",
			Type:        graphql.ID,
		},
	},
})
