package mn_mutation

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_type"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var CreateNote = &graphql.Field{
	Name:        "CreateNote",
	Description: "Create a note. authorize required.",
	Args: graphql.FieldConfigArgument{
		"cw": &graphql.ArgumentConfig{
			Type:         graphql.String,
			DefaultValue: "",
			Description:  "Content Warning. optional.",
		},
		"text": &graphql.ArgumentConfig{
			Type:         graphql.String,
			DefaultValue: "",
			Description:  "Body content.",
		},
		"replyId": &graphql.ArgumentConfig{
			Type:         graphql.ID,
			DefaultValue: "",
			Description:  "Input reply note id. optional.",
		},
		"reTextId": &graphql.ArgumentConfig{
			Type:         graphql.ID,
			DefaultValue: "",
			Description:  "Input reply note id. optional.",
		},
		"localOnly": &graphql.ArgumentConfig{
			Type:         graphql.Boolean,
			DefaultValue: false,
			Description:  "If it is true, don't send fediverse.",
		},
	},
	Type: mn_type.NoteQLType,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, errors.New("authorize required")
		}

		var (
			replyId  = uuid.Nil
			reTextId = uuid.Nil
		)

		replyId, _ = uuid.Parse(p.Args["replyId"].(string))
		reTextId, _ = uuid.Parse(p.Args["reTextId"].(string))

		newNote, err := internal.CreateNoteFromLocal(
			currentUser,
			p.Args["cw"].(string),
			p.Args["text"].(string),
			replyId,
			reTextId,
			p.Args["localOnly"].(bool),
		)
		if err != nil {
			return nil, err
		}

		return mn_type.ConvNoteInternal2GQLType(newNote), nil
	},
}
