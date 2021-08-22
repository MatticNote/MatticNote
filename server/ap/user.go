package ap

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func apUserController(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}
	if misc.IsAPAcceptHeader(c) {
		return RenderUser(c, targetUuid)
	} else {
		targetUser, err := internal.GetLocalUser(targetUuid)
		if err != nil {
			switch err {
			case internal.ErrNoSuchUser:
				return fiber.ErrNotFound
			case internal.ErrUserGone:
				return fiber.ErrGone
			case internal.ErrUserSuspended:
				return fiber.ErrForbidden
			default:
				return err
			}
		}
		return c.Redirect(fmt.Sprintf("/@%s", targetUser.Username))
	}
}

func RenderUser(c *fiber.Ctx, targetUuid uuid.UUID) error {
	targetUser, err := internal.GetLocalUser(targetUuid)
	if err != nil {
		switch err {
		case internal.ErrNoSuchUser:
			return fiber.ErrNotFound
		case internal.ErrUserGone:
			return fiber.ErrGone
		case internal.ErrUserSuspended:
			return fiber.ErrForbidden
		default:
			return err
		}
	}
	targetUserPublicKey, err := internal.GetUserPublicKey(targetUuid)
	if err != nil {
		panic(err)
	}

	c.Set("Content-Type", "application/activity+json; charset=utf-8")

	baseUrl := fmt.Sprintf("%s/activity/user/%s", config.Config.Server.Endpoint, targetUuid.String())

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
			"publicKeyPem": string(pem.EncodeToMemory(targetUserPublicKey)),
		},
	}

	body, err := json.Marshal(renderMap)
	if err != nil {
		return err
	}

	return c.Send(body)
}
