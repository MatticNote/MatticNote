package mn_type

import "github.com/graphql-go/graphql"

type MetaType struct {
	Version string `json:"version"`
}

var MetaQLType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Meta",
	Description: "Meta type",
	Fields: graphql.Fields{
		"version": &graphql.Field{
			Name:        "Version",
			Description: "MatticNote's version",
			Type:        graphql.NewNonNull(graphql.String),
		},
	},
})
