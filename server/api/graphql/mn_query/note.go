package mn_query

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/api/graphql/common"
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_type"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var GetNote = &graphql.Field{
	Name:        "Get note",
	Description: "Get the note.",
	Args: graphql.FieldConfigArgument{
		"noteID": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target noteID",
		},
	},
	Type: graphql.NewNonNull(mn_type.NoteQLType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		targetNoteId, err := uuid.Parse(p.Args["noteID"].(string))
		if err != nil {
			return nil, common.ErrInvalidUUID
		}

		targetNote, err := internal.GetNote(targetNoteId)
		if err != nil {
			return nil, err
		}

		return mn_type.ConvNoteInternal2GQLType(targetNote), nil
	},
}
