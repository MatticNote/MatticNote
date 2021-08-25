package mn_query

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal/user"
	"github.com/MatticNote/MatticNote/server/api/graphql/common"
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

		targetUser, err := user.GetUser(currentUser)
		if err != nil {
			return nil, err
		}

		return mn_type.ConvUsrInternal2GQLType(targetUser), nil
	},
}

var GetUser = &graphql.Field{
	Name:        "Get user",
	Description: "Get the user.",
	Args: graphql.FieldConfigArgument{
		"userID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target userID",
		},
	},
	Type: mn_type.UserQLType,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		targetUserId, err := uuid.Parse(p.Args["userID"].(string))
		if err != nil {
			return nil, common.ErrInvalidUUID
		}

		targetUser, err := user.GetUser(targetUserId)
		if err != nil {
			return nil, err
		}

		return mn_type.ConvUsrInternal2GQLType(targetUser), nil
	},
}
