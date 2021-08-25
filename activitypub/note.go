package activitypub

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
)

const publicUrl = "https://www.w3.org/ns/activitystreams#Public"

var noteContext = []interface{}{
	"https://www.w3.org/ns/activitystreams",
	map[string]interface{}{
		"sensitive":   "as:sensitive",
		"toot":        "http://joinmastodon.org/ns#",
		"votersCount": "toot:votersCount",
	},
}

func parseSender(targetNote *internal.NoteStruct) (to, cc []interface{}) {
	authorBaseUrl := fmt.Sprintf("%s/activity/user/%s", config.Config.Server.Endpoint, targetNote.Author.Uuid.String())

	switch targetNote.Visibility {
	case "FOLLOWER":
		to = []interface{}{
			fmt.Sprintf("%s/followers", authorBaseUrl),
		}
		cc = []interface{}{}
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

	if targetNote.ReText != nil {
		authorBaseUrl := fmt.Sprintf("%s/activity/user/%s", config.Config.Server.Endpoint, targetNote.ReText.Author.Uuid.String())
		cc = append(cc, authorBaseUrl)
	}

	return
}

func RenderNote(targetNote *internal.NoteStruct) map[string]interface{} {
	noteBaseUrl := fmt.Sprintf("%s/activity/note/%s", config.Config.Server.Endpoint, targetNote.Uuid.String())
	authorBaseUrl := fmt.Sprintf("%s/activity/user/%s", config.Config.Server.Endpoint, targetNote.Author.Uuid.String())

	var inReplyUrl *string = nil
	if targetNote.Reply != nil {
		inReplyUrlBase := fmt.Sprintf("%s/activity/note/%s", config.Config.Server.Endpoint, targetNote.Reply.Uuid.String())
		inReplyUrl = &inReplyUrlBase
	}

	to, cc := parseSender(targetNote)

	renderMap := map[string]interface{}{
		"@context":     noteContext,
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

func RenderNoteActivity(targetNote *internal.NoteStruct) map[string]interface{} {
	var renderMap map[string]interface{}
	activityBaseUrl := fmt.Sprintf("%s/activity/note/%s/activity", config.Config.Server.Endpoint, targetNote.Uuid.String())
	authorBaseUrl := fmt.Sprintf("%s/activity/user/%s", config.Config.Server.Endpoint, targetNote.Author.Uuid.String())
	to, cc := parseSender(targetNote)

	if targetNote.ReText == nil {
		object := RenderNote(targetNote)
		delete(object, "@context")

		renderMap = map[string]interface{}{
			"@context":  noteContext,
			"id":        activityBaseUrl,
			"type":      "Create",
			"actor":     authorBaseUrl,
			"published": targetNote.CreatedAt,
			"to":        to,
			"cc":        cc,
			"object":    object,
		}
	} else {
		reTextBaseUrl := fmt.Sprintf("%s/activity/note/%s", config.Config.Server.Endpoint, targetNote.ReText.Uuid)
		renderMap = map[string]interface{}{
			"@context":  noteContext,
			"id":        activityBaseUrl,
			"type":      "Announce",
			"actor":     authorBaseUrl,
			"published": targetNote.CreatedAt,
			"to":        to,
			"cc":        cc,
			"object":    reTextBaseUrl,
		}
	}

	return renderMap
}
