package mn_mutation

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_type"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var FollowUser = &graphql.Field{
	Name:        "FollowUser",
	Description: "Follow the user. authorize required.",
	Args: graphql.FieldConfigArgument{
		"targetID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target userID.",
		},
	},
	Type: graphql.NewNonNull(mn_type.UserCreateFollowRelationQLType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, ErrAuthorizeRequired
		}

		targetID, err := uuid.Parse(p.Args["targetID"].(string))
		if err != nil {
			return nil, errors.New("invalid UUID")
		}

		targetUser, err := internal.GetUser(targetID)
		if err != nil {
			return nil, errors.New("user not found")
		}

		err = internal.CreateFollowRelation(currentUser, targetUser.Uuid, targetUser.AcceptManually)
		if err != nil {
			return nil, err
		}

		result := mn_type.UserCreateFollowRelationType{IsPending: targetUser.AcceptManually}

		return result, nil
	},
}
