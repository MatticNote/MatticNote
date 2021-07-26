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

		return mn_type.ConvUsrInternal2GQLType(targetUser), nil
	},
}

var GetUser = &graphql.Field{
	Name:        "Get the user",
	Description: "Get the user.",
	Args: graphql.FieldConfigArgument{
		"targetID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target userID",
		},
	},
	Type: mn_type.UserQLType,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		targetUserId, err := uuid.Parse(p.Args["targetID"].(string))
		if err != nil {
			return nil, errors.New("invalid uuid")
		}

		targetUser, err := internal.GetUser(targetUserId)
		if err != nil {
			return nil, err
		}

		return mn_type.ConvUsrInternal2GQLType(targetUser), nil
	},
}
