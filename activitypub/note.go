package activitypub

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
)

const publicUrl = "https://www.w3.org/ns/activitystreams#Public"

func RenderNote(targetNote *internal.NoteStruct) map[string]interface{} {
	var (
		to []interface{}
		cc []interface{}
	)

	noteBaseUrl := fmt.Sprintf("%s/activity/note/%s", config.Config.Server.Endpoint, targetNote.Uuid.String())
	authorBaseUrl := fmt.Sprintf("%s/activity/user/%s", config.Config.Server.Endpoint, targetNote.Author.Uuid.String())

	var inReplyUrl *string = nil
	if targetNote.Reply != nil {
		inReplyUrlBase := fmt.Sprintf("%s/activity/note/%s", config.Config.Server.Endpoint, targetNote.Reply.Uuid.String())
		inReplyUrl = &inReplyUrlBase
	}

	switch targetNote.Visibility {
	case "UNLISTED":
		to = []interface{}{
			fmt.Sprintf("%s/followers", authorBaseUrl),
		}
		cc = []interface{}{
			publicUrl,
		}
	case "PUBLIC":
		fallthrough
	default:
		to = []interface{}{
			publicUrl,
		}
		cc = []interface{}{
			fmt.Sprintf("%s/followers", authorBaseUrl),
		}
	}

	renderMap := map[string]interface{}{
		"@context": []interface{}{
			"https://www.w3.org/ns/activitystreams",
		},
		"id":           noteBaseUrl,
		"type":         "Note", // todo: 投票の投稿は`Question`になるので将来的に対応できる仕組みにする
		"url":          fmt.Sprintf("%s/@%s/%s", config.Config.Server.Endpoint, targetNote.Author.Username, targetNote.Uuid.String()),
		"attributedTo": authorBaseUrl,
		"inReplyTo":    inReplyUrl,
		"summary":      targetNote.Cw,
		"content":      targetNote.Body,
		"published":    targetNote.CreatedAt,
		"to":           to,
		"cc":           cc,
	}

	return renderMap
}
