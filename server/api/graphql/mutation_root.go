package graphql

import (
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_mutation"
	"github.com/graphql-go/graphql"
)

var mutationRoot = graphql.ObjectConfig{
	Name:        "MNMutation",
	Description: "MatticNote Mutation",
	Fields: graphql.Fields{
		"createApp":    mn_mutation.CreateApp,
		"createNote":   mn_mutation.CreateNote,
		"deleteNote":   mn_mutation.DeleteNote,
		"followUser":   mn_mutation.FollowUser,
		"unfollowUser": mn_mutation.UnFollowUser,
		"muteUser":     mn_mutation.MuteUser,
		"unmuteUser":   mn_mutation.UnMuteUser,
		"blockUser":    mn_mutation.BlockUser,
		"unblockUser":  mn_mutation.UnBlockUser,
	},
}
