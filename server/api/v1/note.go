package v1

import (
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/MatticNote/MatticNote/internal/note"
	"github.com/MatticNote/MatticNote/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
	"time"
)

type apiV1NoteStruct struct {
	ID        string           `json:"id"`
	Owner     *apiV1UserStruct `json:"owner"`
	CW        *string          `json:"cw"`
	Body      *string          `json:"body"`
	CreatedAt *time.Time       `json:"created_at"`
}

type apiV1NoteCreateForm struct {
	CW   *string
	Body *string
}

func newApiV1NoteStructFromInternal(itn *types.Note, itu *types.User) *apiV1NoteStruct {
	ns := new(apiV1NoteStruct)
	us := newApiV1UserStructFromInternal(itu)

	ns.Owner = us
	ns.ID = itn.ID.String()

	if itn.CW.Valid {
		ns.CW = &itn.CW.String
	}

	if itn.Body.Valid {
		ns.Body = &itn.Body.String
	}

	if itn.CreatedAt.Valid {
		ns.CreatedAt = &itn.CreatedAt.Time
	}

	return ns
}

func noteApiRoute(r fiber.Router) {
	r.Post("/", loginRequired, noteCreate)
	r.Get("/:id", noteGet)
}

func noteCreate(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*types.User)

	form := new(apiV1NoteCreateForm)

	if err := c.BodyParser(form); err != nil {
		return apiBadRequest(c, "Bad form")
	}

	err := validator.New().Struct(form)
	if err != nil {
		return apiBadRequest(c, "Bad form")
	}

	createdNote, err := note.CreateNote(user.ID, form.CW, form.Body)
	if err != nil {
		return err
	}

	noteOwner, err := account.GetUser(createdNote.Owner)
	if err != nil {
		return err
	}

	return c.JSON(newApiV1NoteStructFromInternal(createdNote, noteOwner))
}

func noteGet(c *fiber.Ctx) error {
	noteId, err := ksuid.Parse(c.Params("id"))
	if err != nil {
		return apiNotFound(c, "Note not found")
	}

	noteDetail, err := note.GetNote(noteId)
	if err != nil {
		if err == note.ErrNoteNotFound {
			return apiNotFound(c, "Note not found")
		} else {
			return err
		}
	}

	owner, err := account.GetUser(noteDetail.Owner)
	if err != nil {
		return err
	}

	return c.JSON(newApiV1NoteStructFromInternal(noteDetail, owner))
}
