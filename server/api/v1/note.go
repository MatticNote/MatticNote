package v1

import (
	"errors"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/MatticNote/MatticNote/internal/note"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
	"time"
)

type (
	apiV1NoteStruct struct {
		ID        ksuid.KSUID      `json:"id"`
		Owner     *apiV1UserStruct `json:"owner"`
		CW        *string          `json:"cw"`
		Body      *string          `json:"body"`
		Reply     *apiV1NoteStruct `json:"reply,omitempty"`
		Retext    *apiV1NoteStruct `json:"retext,omitempty"`
		CreatedAt time.Time        `json:"created_at"`
	}

	apiV1NoteCreateStruct struct {
		CW       *string `json:"cw"`
		Body     *string `json:"body"`
		ReplyID  *string `json:"reply_id"`
		RetextID *string `json:"retext_id"`
	}
)

func newApiV1NoteStructFromSchema(it *schemas.Note) *apiV1NoteStruct {
	n := new(apiV1NoteStruct)

	n.ID = it.ID
	n.CreatedAt = it.CreatedAt

	if it.CW.Valid {
		n.CW = &it.CW.String
	}

	if it.Body.Valid {
		n.Body = &it.Body.String
	}

	return n
}

func noteApiRoute(r fiber.Router) {
	r.Post("/", loginRequired, activeAccountRequired, createNote)
	r.Get("/:id", getNote)
}

func createNote(c *fiber.Ctx) error {
	form := new(apiV1NoteCreateStruct)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	err := validator.New().Struct(*form)
	if err != nil {
		return apiBadRequest(c, "Invalid form")
	}

	user := c.Locals("currentUser").(*schemas.User)

	var (
		replyId  *ksuid.KSUID = nil
		retextId *ksuid.KSUID = nil
	)

	if form.ReplyID != nil {
		parse, err := ksuid.Parse(*form.ReplyID)
		if err != nil {
			return apiBadRequest(c, "Invalid reply_id")
		}
		replyId = &parse
	}

	if form.RetextID != nil {
		parse, err := ksuid.Parse(*form.RetextID)
		if err != nil {
			return apiBadRequest(c, "Invalid retext_id")
		}
		retextId = &parse
	}

	newNote, err := note.CreateNote(
		user.ID,
		form.CW,
		form.Body,
		replyId,
		retextId,
	)
	if err != nil {
		switch {
		case errors.Is(err, note.ErrNoteNotFound):
			return apiBadRequest(c, "Note not found for reply_id or retext_id")
		case errors.Is(err, note.ErrNoteInvalidParam):
			return apiBadRequest(c, "Incorrect Note form")
		default:
			return err
		}
	}

	newNoteResponse := newApiV1NoteStructFromSchema(newNote)

	if newNote.Owner != nil {
		noteOwner, err := account.GetUser(*newNote.Owner)
		if err != nil {
			return err
		}
		newNoteResponse.Owner = newApiV1UserStructFromSchema(noteOwner)
	}

	if newNote.ReplyID != nil {
		noteDetailReply, err := note.GetNote(*newNote.ReplyID)
		if err != nil {
			return err
		}
		newNoteResponse.Reply = newApiV1NoteStructFromSchema(noteDetailReply)
	}

	if newNote.RetextID != nil {
		noteDetailRetext, err := note.GetNote(*newNote.RetextID)
		if err != nil {
			return err
		}
		newNoteResponse.Retext = newApiV1NoteStructFromSchema(noteDetailRetext)
	}

	return c.JSON(newNoteResponse)
}

func getNote(c *fiber.Ctx) error {
	id, err := ksuid.Parse(c.Params("id"))
	if err != nil {
		return apiNotFound(c, "Note not found")
	}

	noteDetail, err := note.GetNote(id)
	if err != nil {
		switch {
		case errors.Is(err, note.ErrNoteNotFound):
			return apiNotFound(c, "Note not found")
		default:
			return err
		}
	}

	noteResponse := newApiV1NoteStructFromSchema(noteDetail)

	if noteDetail.Owner != nil {
		noteOwner, err := account.GetUser(*noteDetail.Owner)
		if err != nil {
			return err
		}
		noteResponse.Owner = newApiV1UserStructFromSchema(noteOwner)
	}

	if noteDetail.ReplyID != nil {
		noteDetailReply, err := note.GetNote(*noteDetail.ReplyID)
		if err != nil {
			return err
		}
		noteResponse.Reply = newApiV1NoteStructFromSchema(noteDetailReply)
	}

	if noteDetail.RetextID != nil {
		noteDetailRetext, err := note.GetNote(*noteDetail.ReplyID)
		if err != nil {
			return err
		}
		noteResponse.Retext = newApiV1NoteStructFromSchema(noteDetailRetext)
	}

	return c.JSON(noteResponse)
}
