package activitypub

import (
	"encoding/pem"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
)

func RenderActor(targetUser *internal.LocalUserStruct) (map[string]interface{}, error) {
	baseUrl := fmt.Sprintf("%s/activity/user/%s", config.Config.Server.Endpoint, targetUser.Uuid.String())

	renderMap := fiber.Map{
		"@context": []interface{}{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
			map[string]interface{}{
				"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
				"toot":                      "http://joinmastodon.org/ns#",
				"featured": map[string]interface{}{
					"@id":   "toot:featured",
					"@type": "@id",
				},
				"alsoKnownAs": map[string]interface{}{
					"@id":   "toot:alsoKnownAs",
					"@type": "@id",
				},
				"movedTo": map[string]interface{}{
					"@id":   "toot:movedTo",
					"@type": "@id",
				},
				"schema":        "http://schema.org#",
				"PropertyValue": "schema:PropertyValue",
				"value":         "schema:value",
				"discoverable":  "toot:discoverable",
				"suspended":     "toot:suspended",
			},
		},
		"id": baseUrl,
		"type": func() string {
			if targetUser.IsBot {
				return "Service"
			} else {
				return "Person"
			}
		}(),
		"preferredUsername":         targetUser.Username,
		"manuallyApprovesFollowers": targetUser.AcceptManually,
		"endpoints": map[string]interface{}{
			"sharedInbox": fmt.Sprintf("%s/activity/inbox", config.Config.Server.Endpoint),
		},
		"publicKey": map[string]interface{}{
			"id":           fmt.Sprintf("%s#main-key", baseUrl),
			"owner":        baseUrl,
			"publicKeyPem": string(pem.EncodeToMemory(targetUser.PublicKey)),
		},
		"name":      targetUser.DisplayName,
		"summary":   targetUser.Summary,
		"inbox":     fmt.Sprintf("%s/inbox", baseUrl),
		"suspended": targetUser.IsSuspend,
	}

	return renderMap, nil
}
