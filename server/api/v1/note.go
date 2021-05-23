package v1

import (
	"context"
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"net/http"
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

	newNoteUuid, err := internal.CreateNoteFromLocal(
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

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"uuid": newNoteUuid,
	})
}

func getNote(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return badRequest(c, "Not valid UUID format")
	}

	var (
		noteRes       v1NoteRes
		noteAuthorRes v1UserRes
		isActive      bool
		isSuspend     bool
	)

	err = database.DBPool.QueryRow(
		context.Background(),
		"select note.uuid, note.created_at, cw, body, local_only,"+
			" u.uuid, username, host, display_name, summary, u.created_at, updated_at, accept_manually, is_bot, is_active, is_suspend "+
			"from note, \"user\" u where u.uuid = note.author and note.uuid = $1;",
		targetUuid.String(),
	).Scan(
		&noteRes.Uuid,
		&noteRes.CreatedAt,
		&noteRes.Cw,
		&noteRes.Body,
		&noteRes.LocalOnly,
		&noteAuthorRes.Uuid,
		&noteAuthorRes.Username,
		&noteAuthorRes.Host,
		&noteAuthorRes.DisplayName,
		&noteAuthorRes.Summary,
		&noteAuthorRes.CreatedAt,
		&noteAuthorRes.UpdatedAt,
		&noteAuthorRes.AcceptManually,
		&noteAuthorRes.IsBot,
		&isActive,
		&isSuspend,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return notFound(c, "Specified user was not found")
		} else {
			return err
		}
	}

	if !isActive {
		c.Status(http.StatusGone)
		return nil
	}

	if isSuspend {
		return forbidden(c, "Note Author is Suspend")
	}

	noteRes.Author = noteAuthorRes

	return c.Status(http.StatusOK).JSON(noteRes)
}
