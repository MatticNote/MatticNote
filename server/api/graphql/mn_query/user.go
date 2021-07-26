package mn_query

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_type"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var CurrentUser = &graphql.Field{
	Name:        "Current User",
	Description: "Get current User ID. Require authentication.",
	Type:        graphql.NewNonNull(mn_type.UserQLType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, errors.New("authorize required")
		}

		targetUser, err := internal.GetUser(currentUser)
		if err != nil {
			return nil, err
		}

		output := mn_type.ConvUsrInternal2GQLType(targetUser)

		return output, nil
	},
}
