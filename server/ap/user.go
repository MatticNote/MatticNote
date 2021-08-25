package ap

import (
	"encoding/json"
	"fmt"
	"github.com/MatticNote/MatticNote/activitypub"
	"github.com/MatticNote/MatticNote/internal/ist"
	"github.com/MatticNote/MatticNote/internal/user"
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
	targetUser, err := user.GetLocalUser(targetUuid)
	if err != nil {
		switch err {
		case user.ErrNoSuchUser:
			return fiber.ErrNotFound
		case user.ErrUserGone:
			return fiber.ErrGone
		default:
			return err
		}
	}
	c.Locals("targetUser", targetUser)
	return c.Next()
}

func apUserController(c *fiber.Ctx) error {
	targetUser := c.Locals("targetUser").(*ist.LocalUserStruct)
	if misc.IsAPAcceptHeader(c) {
		return RenderUser(c, targetUser)
	} else {
		return c.Redirect(fmt.Sprintf("/@%s", targetUser.Username))
	}
}

func RenderUser(c *fiber.Ctx, targetUser *ist.LocalUserStruct) error {
	c.Set("Content-Type", "application/activity+json; charset=utf-8")

	actor := activitypub.RenderActor(targetUser)

	body, err := json.Marshal(actor)
	if err != nil {
		return err
	}

	return c.Send(body)
}
