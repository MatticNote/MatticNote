package v1

import (
	"fmt"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
)

type newNoteStruct struct {
	Cw         string
	Text       string `validate:"required"`
	ReplyUuid  string
	ReTextUuid string
	LocalOnly  bool
}

func postNote(c *fiber.Ctx) error {
	currentUsr, ok := c.Locals(loginUserLocal).(*internal.LocalUserStruct)
	if !ok {
		return unauthorized(c)
	}
	formData := new(newNoteStruct)

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
