package v1

import (
	"fmt"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func postNote(c *fiber.Ctx) error {
	currentUsr, ok := c.Locals(internal.LoginUserLocal).(*internal.LocalUserStruct)
	if !ok {
		return unauthorized(c)
	}
	formData := new(newNoteReq)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(*formData); errs != nil {
		return badRequest(c, fmt.Sprintf("invalid form: %v", errs))
	}

	replyUuid, err := uuid.Parse(formData.ReplyUuid)
	if err != nil {
		replyUuid = uuid.Nil
	}

	reTextUuid, err := uuid.Parse(formData.ReTextUuid)
	if err != nil {
		reTextUuid = uuid.Nil
	}

	newNote, err := internal.CreateNoteFromLocal(
		currentUsr.Uuid,
		formData.Cw,
		formData.Text,
		replyUuid,
		reTextUuid,
		formData.LocalOnly,
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(convFromInternalNote(*newNote))
}

func getNote(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return badRequest(c, "Not valid UUID format")
	}

	noteRes, err := internal.GetNote(targetUuid)
	if err != nil {
		switch err {
		case internal.ErrNoteNotFound:
			return notFound(c, "no such note")
		case internal.ErrUserGone:
			c.Status(fiber.StatusGone)
			return nil
		case internal.ErrUserSuspended:
			return forbidden(c, "Note author is suspended")
		default:
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(noteRes)
}

func deleteNote(c *fiber.Ctx) error {
	currentUsr, ok := c.Locals(internal.LoginUserLocal).(*internal.LocalUserStruct)
	if !ok {
		return unauthorized(c)
	}

	targetNoteUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return badRequest(c, "Not valid UUID format")
	}

	noteRes, err := internal.GetNote(targetNoteUuid)
	if err != nil {
		switch err {
		case internal.ErrNoteNotFound:
			return notFound(c, "no such note")
		case internal.ErrUserGone:
			c.Status(fiber.StatusGone)
			return nil
		case internal.ErrUserSuspended:
			return forbidden(c, "Note author is suspended")
		default:
			return nil
		}
	}

	if currentUsr.Uuid != noteRes.Author.Uuid {
		return forbidden(c, "this note is not owner")
	}

	err = internal.DeleteNote(noteRes.Uuid)
	if err != nil {
		return err
	}

	c.Status(fiber.StatusNoContent)
	return nil
}
