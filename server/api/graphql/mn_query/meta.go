package mn_query

import (
	"fmt"
	"github.com/MatticNote/MatticNote/internal/version"
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_type"
	"github.com/graphql-go/graphql"
)

var Meta = &graphql.Field{
	Name:        "Meta",
	Description: "Meta information",
	Type:        graphql.NewNonNull(mn_type.MetaQLType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return mn_type.MetaType{
			Version: fmt.Sprintf("%s-%s", version.Version, version.Revision),
		}, nil
	},
}
