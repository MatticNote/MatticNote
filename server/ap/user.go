package ap

import (
	"encoding/json"
	"fmt"
	"github.com/MatticNote/MatticNote/activitypub"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func apUserHandler(c *fiber.Ctx) error {
	c.Query("uuid")
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}
	targetUser, err := internal.GetLocalUser(targetUuid)
	if err != nil {
		switch err {
		case internal.ErrNoSuchUser:
			return fiber.ErrNotFound
		case internal.ErrUserGone:
			return fiber.ErrGone
		default:
			return err
		}
	}
	c.Locals("targetUser", targetUser)
	return c.Next()
}

func apUserController(c *fiber.Ctx) error {
	targetUser := c.Locals("targetUser").(*internal.LocalUserStruct)
	if misc.IsAPAcceptHeader(c) {
		return RenderUser(c, targetUser)
	} else {
		return c.Redirect(fmt.Sprintf("/@%s", targetUser.Username))
	}
}

func RenderUser(c *fiber.Ctx, targetUser *internal.LocalUserStruct) error {
	c.Set("Content-Type", "application/activity+json; charset=utf-8")

	actor, err := activitypub.RenderActor(targetUser)
	if err != nil {
		return err
	}

	body, err := json.Marshal(actor)
	if err != nil {
		return err
	}

	return c.Send(body)
}
