package server

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/MatticNote/MatticNote/server/ap"
	"github.com/gofiber/fiber/v2"
)

func userProfileController(c *fiber.Ctx) error {
	if misc.IsAPAcceptHeader(c) {
		// ActivityPub Render
		targetUser, err := internal.GetLocalUserFromUsername(c.Params("username"))
		if err != nil && err == internal.ErrNoSuchUser {
			return fiber.ErrNotFound
		} else if err != nil {
			return err
		}
		return ap.RenderUser(c, targetUser)
	} else {
		// Normal render
		return userProfileView(c)
	}
}

func userProfileView(c *fiber.Ctx) error {
	targetUser, err := internal.GetLocalUserFromUsername(c.Params("username"))
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

	return c.Render(
		"user_profile",
		fiber.Map{
			"username":       targetUser.Username,
			"displayName":    targetUser.DisplayName.String,
			"summary":        targetUser.Summary.String,
			"createdAt":      targetUser.CreatedAt.Time,
			"updatedAt":      targetUser.UpdatedAt.Time,
			"isSuperUser":    targetUser.IsSuperuser,
			"isBot":          targetUser.IsBot,
			"acceptManually": targetUser.AcceptManually,
		},
	)
}
