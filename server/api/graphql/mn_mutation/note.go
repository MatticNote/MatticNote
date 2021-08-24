package mn_mutation

import (
	"errors"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/api/graphql/common"
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
		"visibility": &graphql.ArgumentConfig{
			Type:         graphql.String,
			DefaultValue: "public",
			Description:  "Visibility",
		},
	},
	Type: graphql.NewNonNull(mn_type.NoteQLType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		currentUser, ok := p.Context.Value("currentUser").(uuid.UUID)
		if !ok {
			return nil, common.ErrAuthorizeRequired
		}

		var (
			replyId  *uuid.UUID
			reTextId *uuid.UUID
			cw       *string
			text     *string
		)

		replyIdRaw, err := uuid.Parse(p.Args["replyId"].(string))
		if err == nil {
			replyId = &replyIdRaw
		}
		reTextIdRaw, err := uuid.Parse(p.Args["reTextId"].(string))
		if err == nil {
			reTextId = &reTextIdRaw
		}
		cwRaw, ok := p.Args["cw"].(string)
		if ok {
			cw = &cwRaw
		}
		textRaw, ok := p.Args["text"].(string)
		if ok {
			text = &textRaw
		}

		newNote, err := internal.CreateNoteFromLocal(
			currentUser,
			cw,
			text,
			replyId,
			reTextId,
			p.Args["localOnly"].(bool),
			p.Args["visibility"].(string),
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
			return nil, common.ErrAuthorizeRequired
		}

		targetNoteUuid, err := uuid.Parse(p.Args["noteId"].(string))
		if err != nil {
			return nil, common.ErrInvalidUUID
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
