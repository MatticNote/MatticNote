package graphql

import (
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_query"
	"github.com/graphql-go/graphql"
)

var queryRoot = graphql.ObjectConfig{
	Name:        "MNQuery",
	Description: "MatticNote Query",
	Fields: graphql.Fields{
		"meta":        mn_query.Meta,
		"currentUser": mn_query.CurrentUser,
		"getUser":     mn_query.GetUser,
		"getNote":     mn_query.GetNote,
	},
}
