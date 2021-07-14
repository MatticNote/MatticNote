package mn_query

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_misc"
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_type"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var CurrentUser = &graphql.Field{
	Name:        "Current User",
	Type:        graphql.NewNonNull(mn_type.UserQLType),
	Description: "Get current User ID. Require authentication.",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, errors.New("authorize required")
		}

		targetUser, err := internal.GetLocalUser(currentUser)
		if err != nil {
			return nil, err
		}

		output := mn_type.UserType{
			Uuid:           targetUser.Uuid,
			Username:       targetUser.Username,
			Host:           mn_misc.Conv2Interface(targetUser.Host),
			DisplayName:    mn_misc.Conv2Interface(targetUser.DisplayName),
			Summary:        mn_misc.Conv2Interface(targetUser.Summary),
			CreatedAt:      mn_misc.Conv2Interface(targetUser.CreatedAt),
			UpdatedAt:      mn_misc.Conv2Interface(targetUser.UpdatedAt),
			AcceptManually: targetUser.AcceptManually,
			IsBot:          targetUser.IsBot,
		}

		return output, nil
	},
}
