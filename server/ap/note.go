package ap

import (
	"encoding/json"
	"fmt"
	"github.com/MatticNote/MatticNote/activitypub"
	"github.com/MatticNote/MatticNote/internal/ist"
	iNote "github.com/MatticNote/MatticNote/internal/note"
	"github.com/MatticNote/MatticNote/internal/user"
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
	targetNote, err := iNote.GetNote(targetUuid, true)
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
	// Ignore remote user's note or local only
	if targetNote.Author.Host.Valid || targetNote.LocalOnly {
		return fiber.ErrNotFound
	}
	c.Locals("targetNote", targetNote)
	return c.Next()
}

func apNoteController(c *fiber.Ctx) error {
	targetNote := c.Locals("targetNote").(*ist.NoteStruct)
	if misc.IsAPAcceptHeader(c) {
		return RenderNote(c, targetNote)
	} else {
		return c.Redirect(fmt.Sprintf("/@%s/%s", targetNote.Author.Username, targetNote.Uuid.String()))
	}
}

func RenderNote(c *fiber.Ctx, targetNote *ist.NoteStruct) error {
	c.Set("Content-Type", "application/activity+json; charset=utf-8")

	if targetNote.ReText != nil {
		return c.Redirect(fmt.Sprintf("/activity/note/%s", targetNote.ReText.Uuid.String()))
	}

	note := activitypub.RenderNote(targetNote)

	body, err := json.Marshal(note)
	if err != nil {
		return err
	}

	return c.Send(body)
}

func apNoteActivityController(c *fiber.Ctx) error {
	targetNote := c.Locals("targetNote").(*ist.NoteStruct)
	return renderNoteActivity(c, targetNote)
}

func renderNoteActivity(c *fiber.Ctx, targetNote *ist.NoteStruct) error {
	c.Set("Content-Type", "application/activity+json; charset=utf-8")

	activity := activitypub.RenderNoteActivity(targetNote)

	body, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	return c.Send(body)
}
