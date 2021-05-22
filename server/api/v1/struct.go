package v1

import (
	"github.com/MatticNote/MatticNote/misc"
)

type (
	v1UserRes struct {
		Uuid           string          `json:"uuid"`
		Username       string          `json:"username"`
		Host           misc.NullString `json:"host"`
		DisplayName    misc.NullString `json:"display_name"`
		Summary        misc.NullString `json:"summary"`
		CreatedAt      misc.NullTime   `json:"created_at"`
		UpdatedAt      misc.NullTime   `json:"updated_at"`
		AcceptManually bool            `json:"accept_manually"`
		IsBot          bool            `json:"is_bot"`
	}

	v1NoteRes struct {
		Uuid      string          `json:"uuid"`
		Author    v1UserRes       `json:"author"`
		CreatedAt misc.NullTime   `json:"created_at"`
		Cw        misc.NullString `json:"cw"`
		Body      misc.NullString `json:"body"`
		LocalOnly bool            `json:"local_only"`
	}
)

type (
	newNoteReq struct {
		Cw         string
		Text       string `validate:"required"`
		ReplyUuid  string `json:"reply_uuid"`
		ReTextUuid string `json:"re_text_uuid"`
		LocalOnly  bool   `json:"local_only"`
	}
)
