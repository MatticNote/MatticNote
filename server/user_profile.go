package server

import (
	"fmt"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/MatticNote/MatticNote/server/ap"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func userProfileHandler(c *fiber.Ctx) error {
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
	c.Locals("targetUser", targetUser)
	return c.Next()
}

func userProfileController(c *fiber.Ctx) error {
	targetUser := c.Locals("targetUser").(*internal.LocalUserStruct)
	if misc.IsAPAcceptHeader(c) {
		return ap.RenderUser(c, targetUser)
	} else {
		return userProfileView(c, targetUser)
	}
}

func userProfileView(c *fiber.Ctx, targetUser *internal.LocalUserStruct) error {
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

func userProfileNoteController(c *fiber.Ctx) error {
	targetUser := c.Locals("targetUser").(*internal.LocalUserStruct)
	noteUuid, err := uuid.Parse(c.Params("noteUuid"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	targetNote, err := internal.GetNote(noteUuid)
	if err != nil {
		switch err {
		case internal.ErrNoteNotFound:
			return fiber.ErrNotFound
		case internal.ErrUserSuspended:
			return fiber.ErrForbidden
		default:
			return err
		}
	}

	if targetUser.Uuid != targetNote.Author.Uuid {
		return c.Redirect(fmt.Sprintf("/@%s/%s", targetNote.Author.Username, targetNote.Uuid))
	}

	if misc.IsAPAcceptHeader(c) {
		return ap.RenderNote(c, targetNote)
	} else {
		return userProfileNoteView(c, targetNote)
	}
}

func userProfileNoteView(c *fiber.Ctx, targetNote *internal.NoteStruct) error {
	return c.Render(
		"user_profile_note",
		fiber.Map{
			"body": targetNote.Body.String,
		},
	)
}
