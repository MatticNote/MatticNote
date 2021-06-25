package v1

import (
	"github.com/MatticNote/MatticNote/internal"
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

	v1UserRelationRes struct {
		Following     bool `json:"following"`
		FollowPending bool `json:"follow_pending"`
		Follows       bool `json:"follows"`
		Muting        bool `json:"muting"`
		Blocking      bool `json:"blocking"`
		Blocked       bool `json:"blocked"`
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

func convFromInternalUser(iu internal.UserStruct) *v1UserRes {
	return &v1UserRes{
		Uuid:           iu.Uuid.String(),
		Username:       iu.Username,
		Host:           iu.Host,
		DisplayName:    iu.DisplayName,
		Summary:        iu.Summary,
		CreatedAt:      iu.CreatedAt,
		UpdatedAt:      iu.UpdatedAt,
		AcceptManually: iu.AcceptManually,
		IsBot:          iu.IsBot,
	}
}

func convFromInternalNote(ns internal.NoteStruct) *v1NoteRes {
	return &v1NoteRes{
		Uuid:      ns.Uuid.String(),
		Author:    *convFromInternalUser(ns.Author),
		CreatedAt: ns.CreatedAt,
		Cw:        ns.Cw,
		Body:      ns.Body,
		LocalOnly: ns.LocalOnly,
	}
}

func convFromInternalUserRelate(ur internal.UserRelationStruct) *v1UserRelationRes {
	return &v1UserRelationRes{
		Following:     ur.Following,
		FollowPending: ur.FollowPending,
		Follows:       ur.Follows,
		Muting:        ur.Muting,
		Blocking:      ur.Blocking,
		Blocked:       ur.Blocked,
	}
}
