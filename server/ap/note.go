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

func apNoteHandler(c *fiber.Ctx) error {
	c.Query("uuid")
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}
	targetNote, err := internal.GetNote(targetUuid)
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
	c.Locals("targetNote", targetNote)
	return c.Next()
}

func apNoteController(c *fiber.Ctx) error {
	targetNote := c.Locals("targetNote").(*internal.NoteStruct)
	if misc.IsAPAcceptHeader(c) {
		return RenderNote(c, targetNote)
	} else {
		return c.Redirect(fmt.Sprintf("/@%s/%s", targetNote.Author.Username, targetNote.Uuid.String()))
	}
}

func RenderNote(c *fiber.Ctx, targetNote *internal.NoteStruct) error {
	c.Set("Content-Type", "application/activity+json; charset=utf-8")

	note := activitypub.RenderNote(targetNote)

	body, err := json.Marshal(note)
	if err != nil {
		return err
	}

	return c.Send(body)
}
