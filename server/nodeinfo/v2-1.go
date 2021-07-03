package nodeinfo

import (
	"encoding/json"
	"fmt"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
)

func nodeinfoV21(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; profile=\"http://nodeinfo.diaspora.software/ns/schema/2.1#\"")
	// source: https://nodeinfo.diaspora.software/docson/index.html#/ns/schema/2.1#$$expand
	niJson, err := json.Marshal(fiber.Map{
		"version": 2.1,
		"software": fiber.Map{
			"name":    "matticnote",
			"version": fmt.Sprintf("%s-%s", internal.Version, internal.Revision),
			// TODO: repository, homepageも設定する
		},
		"protocols": []string{
			"activitypub",
		},
		"services": fiber.Map{
			"inbound":  make([]string, 0),
			"outbound": make([]string, 0),
		},
		"openRegistrations": true, // TODO: 設定ファイルから読み取れるようにする
		"usage": fiber.Map{
			"users": fiber.Map{
				"total": 0, // TODO: ユーザーアカウント数の合計値を関数とかで読めるようにする
				// TODO: activeHalfyearやactiveMonthもできれば
			},
			"localPosts": 0, // TODO: ノート(ローカル)の総数表示
		},
		"metadata": fiber.Map{
			// TODO: MatticNoteのカスタムの項目を追加する
		},
	})
	if err != nil {
		return err
	}

	return c.Send(niJson)
}
