package mn_query

import (
	"github.com/MatticNote/MatticNote/internal/note"
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
		if _, ok := p.Context.Value(common.ContextCurrentUser).(uuid.UUID); ok {
			if err := common.ScopeCheck(p, "read.note"); err != nil {
				return nil, err
			}
		}
		targetNoteId, err := uuid.Parse(p.Args["noteID"].(string))
		if err != nil {
			return nil, common.ErrInvalidUUID
		}

		targetNote, err := note.GetNote(targetNoteId)
		if err != nil {
			return nil, err
		}

		return mn_type.ConvNoteInternal2GQLType(targetNote), nil
	},
}
