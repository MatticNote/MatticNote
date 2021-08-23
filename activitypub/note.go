package activitypub

import (
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
)

func RenderNote(targetNote *internal.NoteStruct) map[string]interface{} {
	noteBaseUrl := fmt.Sprintf("%s/activity/note/%s", config.Config.Server.Endpoint, targetNote.Uuid.String())
	authorBaseUrl := fmt.Sprintf("%s/activity/user/%s", config.Config.Server.Endpoint, targetNote.Author.Uuid.String())

	renderMap := map[string]interface{}{
		"@context": []interface{}{
			"https://www.w3.org/ns/activitystreams",
		},
		"id":           noteBaseUrl,
		"type":         "Note", // todo: 投票の投稿は`Question`になるので将来的に対応できる仕組みにする
		"url":          fmt.Sprintf("%s/@%s/%s", config.Config.Server.Endpoint, targetNote.Author.Username, targetNote.Uuid.String()),
		"attributedTo": authorBaseUrl,
		"summary":      targetNote.Cw,
		"content":      targetNote.Body,
		"published":    targetNote.CreatedAt,
		// todo: 投稿範囲に応じてtoやccを変える
		"to": []interface{}{
			"https://www.w3.org/ns/activitystreams#Public",
		},
		"cc": []interface{}{},
	}

	return renderMap
}
