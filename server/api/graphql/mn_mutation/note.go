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
	Type: graphql.NewNonNull(mn_type.NoteQLType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, ErrAuthorizeRequired
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

var DeleteNote = &graphql.Field{
	Name:        "DeleteNote",
	Description: "Delete a note. authorize required.",
	Args: graphql.FieldConfigArgument{
		"noteId": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "Target note UUID",
		},
	},
	Type: graphql.NewNonNull(graphql.Boolean),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, ErrAuthorizeRequired
		}

		targetNoteUuid, err := uuid.Parse(p.Args["noteId"].(string))
		if err != nil {
			return nil, errors.New("invalid UUID")
		}

		targetNote, err := internal.GetNote(targetNoteUuid)
		if err != nil {
			return nil, err
		}
		if currentUser != targetNote.Author.Uuid {
			return nil, errors.New("specified note is not owned")
		}

		err = internal.DeleteNote(targetNote.Uuid)
		if err != nil {
			return nil, err
		}

		return true, err
	},
}
