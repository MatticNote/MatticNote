package server

import (
	"errors"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/MatticNote/MatticNote/internal/note"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
)

func getUserProfile(c *fiber.Ctx) error {
	user, err := account.GetUserByUsername(c.Params("username"))
	if err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound):
			return fiber.ErrNotFound
		default:
			return err
		}
	}

	return c.Render(
		"user-profile",
		fiber.Map{
			"user": user,
			"ui":   UIFileManifest,
		},
	)
}

func getUserNote(c *fiber.Ctx) error {
	user, err := account.GetUserByUsername(c.Params("username"))
	if err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound):
			return fiber.ErrNotFound
		default:
			return err
		}
	}

	noteId, err := ksuid.Parse(c.Params("noteId"))
	if err != nil {
		return fiber.ErrNotFound
	}

	noteData, err := note.GetNote(noteId)
	if err != nil {
		switch {
		case errors.Is(err, note.ErrNoteNotFound):
			return fiber.ErrNotFound
		default:
			return err
		}
	}

	if noteData.Owner == nil || (noteData.Owner != nil && user.ID != *noteData.Owner) {
		return fiber.ErrNotFound
	}

	var (
		replyNoteData  *schemas.Note
		retextNoteData *schemas.Note
	)

	if noteData.ReplyID != nil {
		replyNoteData, err = note.GetNote(*noteData.ReplyID)
		if err != nil {
			return err
		}
	}

	if noteData.RetextID != nil {
		retextNoteData, err = note.GetNote(*noteData.RetextID)
		if err != nil {
			return err
		}
	}

	return c.Render(
		"user-note",
		fiber.Map{
			"user":       user,
			"note":       noteData,
			"replyNote":  replyNoteData,
			"retextNote": retextNoteData,
			"ui":         UIFileManifest,
		},
	)
}
