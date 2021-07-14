package mn_type

import (
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

type UserType struct {
	Uuid           uuid.UUID   `json:"uuid"`
	Username       string      `json:"username"`
	Host           interface{} `json:"host"`
	DisplayName    interface{} `json:"displayName"`
	Summary        interface{} `json:"summary"`
	CreatedAt      interface{} `json:"createdAt"`
	UpdatedAt      interface{} `json:"updatedAt"`
	AcceptManually bool        `json:"acceptManually"`
	IsBot          bool        `json:"isBot"`
}

var UserQLType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "User type",
	Fields: graphql.Fields{
		"uuid": &graphql.Field{
			Name:        "UUID",
			Description: "User's UUID",
			Type:        graphql.NewNonNull(graphql.ID),
		},
		"username": &graphql.Field{
			Name:        "Username",
			Description: "User's username. not include hostname.",
			Type:        graphql.NewNonNull(graphql.String),
		},
		"host": &graphql.Field{
			Name:        "Host",
			Description: "User's hostname. blank is local user.",
			Type:        graphql.String,
		},
		"displayName": &graphql.Field{
			Name:        "Display Name",
			Description: "Display name.",
			Type:        graphql.String,
		},
		"createdAt": &graphql.Field{
			Name:        "Created at",
			Description: "Created datetime",
			Type:        graphql.DateTime,
		},
		"updatedAt": &graphql.Field{
			Name:        "Updated at",
			Description: "Updated datetime",
			Type:        graphql.DateTime,
		},
		"acceptManually": &graphql.Field{
			Name:        "Accept manually",
			Description: "If it is true, required follow approve from target user.",
			Type:        graphql.NewNonNull(graphql.Boolean),
		},
		"isBot": &graphql.Field{
			Name:        "Is bot",
			Description: "Does this account mainly use a system that posts automatically",
			Type:        graphql.NewNonNull(graphql.Boolean),
		},
	},
})
