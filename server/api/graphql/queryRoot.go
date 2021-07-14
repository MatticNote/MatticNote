package graphql

import (
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_query"
	"github.com/graphql-go/graphql"
)

var queryRoot = graphql.ObjectConfig{
	Name: "MNQuery",
	Fields: graphql.Fields{
		"currentUser": mn_query.CurrentUser,
	},
	Description: "MatticNote Query",
}
