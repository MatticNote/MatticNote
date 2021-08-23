package ap

import (
	"github.com/MatticNote/MatticNote/internal"
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
	return nil
}
