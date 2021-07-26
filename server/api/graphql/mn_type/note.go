package mn_type

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_misc"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

type NoteType struct {
	Uuid      uuid.UUID   `json:"uuid"`
	CreatedAt interface{} `json:"created_at"`
	Cw        interface{} `json:"cw"`
	Text      interface{} `json:"text"`
	ReplyId   *uuid.UUID  `json:"reply_id"`
	ReTextId  *uuid.UUID  `json:"retext_id"`
	LocalOnly bool        `json:"local_only"`
}

var NoteQLType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Note",
	Description: "Note type",
	Fields: graphql.Fields{
		"uuid": &graphql.Field{
			Name:        "UUID",
			Description: "Note UUID",
			Type:        graphql.NewNonNull(graphql.ID),
		},
		"createdAt": &graphql.Field{
			Name:        "Created At",
			Description: "Note created time",
			Type:        graphql.DateTime,
		},
		"cw": &graphql.Field{
			Name:        "CW",
			Description: "Content Warning",
			Type:        graphql.String,
		},
		"text": &graphql.Field{
			Name:        "Text",
			Description: "body",
			Type:        graphql.String,
		},
		"replyId": &graphql.Field{
			Name:        "Reply ID",
			Description: "Reply note ID",
			Type:        graphql.ID,
		},
		"reTextId": &graphql.Field{
			Name:        "ReText ID",
			Description: "ReText note ID",
			Type:        graphql.ID,
		},
		"localOnly": &graphql.Field{
			Name:        "Local Only",
			Description: "If it is true, don't send fediverse.",
			Type:        graphql.NewNonNull(graphql.Boolean),
		},
	},
})

func ConvNoteInternal2GQLType(ins *internal.NoteStruct) NoteType {
	return NoteType{
		Uuid:      ins.Uuid,
		CreatedAt: mn_misc.Conv2Interface(ins.CreatedAt),
		Cw:        mn_misc.Conv2Interface(ins.Cw),
		Text:      mn_misc.Conv2Interface(ins.Body),
		ReplyId:   nil,
		ReTextId:  nil,
		LocalOnly: ins.LocalOnly,
	}
}
