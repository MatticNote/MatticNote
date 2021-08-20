package mn_mutation

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/api/graphql/common"
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
			return nil, common.ErrAuthorizeRequired
		}

		targetID, err := uuid.Parse(p.Args["targetID"].(string))
		if err != nil {
			return nil, common.ErrInvalidUUID
		}

		targetUser, err := internal.GetUser(targetID)
		if err != nil {
			return nil, common.ErrUserNotFound
		}

		err = internal.CreateFollowRelation(currentUser, targetUser.Uuid, targetUser.AcceptManually)
		if err != nil {
			return nil, err
		}

		result := mn_type.UserCreateFollowRelationType{IsPending: targetUser.AcceptManually}

		return result, nil
	},
}

var UnFollowUser = &graphql.Field{
	Name:        "UnFollowUser",
	Description: "Unfollow the user. authorize required. Always returns true on success",
	Args: graphql.FieldConfigArgument{
		"targetID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target userID.",
		},
	},
	Type: graphql.NewNonNull(graphql.Boolean),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, common.ErrAuthorizeRequired
		}

		targetID, err := uuid.Parse(p.Args["targetID"].(string))
		if err != nil {
			return nil, common.ErrInvalidUUID
		}

		targetUser, err := internal.GetUser(targetID)
		if err != nil {
			return nil, common.ErrUserNotFound
		}

		err = internal.DestroyFollowRelation(currentUser, targetUser.Uuid)
		if err != nil {
			return nil, err
		}

		return true, nil
	},
}

var MuteUser = &graphql.Field{
	Name:        "MuteUser",
	Description: "Mute the user. authorize required. Always returns true on success",
	Args: graphql.FieldConfigArgument{
		"targetID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target userID.",
		},
	},
	Type: graphql.NewNonNull(graphql.Boolean),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, common.ErrAuthorizeRequired
		}

		targetID, err := uuid.Parse(p.Args["targetID"].(string))
		if err != nil {
			return nil, common.ErrInvalidUUID
		}

		targetUser, err := internal.GetUser(targetID)
		if err != nil {
			return nil, common.ErrUserNotFound
		}

		err = internal.CreateMuteRelation(currentUser, targetUser.Uuid)
		if err != nil {
			return nil, err
		}

		return true, nil
	},
}

var UnMuteUser = &graphql.Field{
	Name:        "UnMuteUser",
	Description: "Unmute the user. authorize required. Always returns true on success",
	Args: graphql.FieldConfigArgument{
		"targetID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target userID.",
		},
	},
	Type: graphql.NewNonNull(graphql.Boolean),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, common.ErrAuthorizeRequired
		}

		targetID, err := uuid.Parse(p.Args["targetID"].(string))
		if err != nil {
			return nil, common.ErrInvalidUUID
		}

		targetUser, err := internal.GetUser(targetID)
		if err != nil {
			return nil, common.ErrUserNotFound
		}

		err = internal.DestroyMuteRelation(currentUser, targetUser.Uuid)
		if err != nil {
			return nil, err
		}

		return true, nil
	},
}

var BlockUser = &graphql.Field{
	Name:        "BlockUser",
	Description: "Block the user. authorize required. Always returns true on success",
	Args: graphql.FieldConfigArgument{
		"targetID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target userID.",
		},
	},
	Type: graphql.NewNonNull(graphql.Boolean),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, common.ErrAuthorizeRequired
		}

		targetID, err := uuid.Parse(p.Args["targetID"].(string))
		if err != nil {
			return nil, common.ErrInvalidUUID
		}

		targetUser, err := internal.GetUser(targetID)
		if err != nil {
			return nil, common.ErrUserNotFound
		}

		err = internal.CreateBlockRelation(currentUser, targetUser.Uuid)
		if err != nil {
			return nil, err
		}

		return true, nil
	},
}

var UnBlockUser = &graphql.Field{
	Name:        "UnBlockUser",
	Description: "Unlock the user. authorize required. Always returns true on success",
	Args: graphql.FieldConfigArgument{
		"targetID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target userID.",
		},
	},
	Type: graphql.NewNonNull(graphql.Boolean),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, common.ErrAuthorizeRequired
		}

		targetID, err := uuid.Parse(p.Args["targetID"].(string))
		if err != nil {
			return nil, common.ErrInvalidUUID
		}

		targetUser, err := internal.GetUser(targetID)
		if err != nil {
			return nil, common.ErrUserNotFound
		}

		err = internal.DestroyMuteRelation(currentUser, targetUser.Uuid)
		if err != nil {
			return nil, err
		}

		return true, nil
	},
}
