package note

import (
	"database/sql"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/lib/pq"
	"github.com/segmentio/ksuid"
)

var (
	ErrNoteNotFound     = errors.New("note not found")
	ErrNoteInvalidParam = errors.New("note invalid parameter")
)

func CreateNote(
	userId ksuid.KSUID,
	cw,
	body *string,
	replyId *ksuid.KSUID,
	retextId *ksuid.KSUID,
) (*schemas.Note, error) {
	if cw == nil && body == nil && replyId == nil && retextId == nil {
		return nil, ErrNoteInvalidParam
	} else if replyId != nil && (body == nil || body != nil && *body == "") {
		return nil, ErrNoteInvalidParam
	} else if replyId != nil && retextId != nil {
		return nil, ErrNoteInvalidParam
	} else if cw != nil && (body == nil || body != nil && *body == "") {
		return nil, ErrNoteInvalidParam
	}

	newNote := new(schemas.Note)
	newNote.ID = ksuid.New()

	err := database.Database.QueryRow(
		"INSERT INTO notes(id, owner, cw, body, reply_id, retext_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at, owner",
		newNote.ID.String(),
		userId,
		cw,
		body,
		replyId,
		retextId,
	).Scan(&newNote.CreatedAt, &newNote.Owner)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Code.Name() {
			case "foreign_key_violation":
				return nil, ErrNoteNotFound
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if cw != nil {
		err = newNote.CW.Scan(*cw)
		if err != nil {
			return nil, err
		}
	}

	if body != nil {
		err = newNote.Body.Scan(*body)
		if err != nil {
			return nil, err
		}
	}

	newNote.ReplyID = replyId
	newNote.RetextID = retextId

	return newNote, nil
}

func GetNote(
	id ksuid.KSUID,
) (*schemas.Note, error) {
	note := new(schemas.Note)

	err := database.Database.QueryRow(
		"SELECT id, owner, cw, body, reply_id, retext_id, created_at FROM notes WHERE id = $1;",
		id.String(),
	).Scan(
		&note.ID,
		&note.Owner,
		&note.CW,
		&note.Body,
		&note.ReplyID,
		&note.RetextID,
		&note.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoteNotFound
		} else {
			return nil, err
		}
	}

	return note, nil
}

func DeleteNote(
	id ksuid.KSUID,
) error {
	exec, err := database.Database.Exec(
		"DELETE FROM notes WHERE id = $1",
		id.String(),
	)
	if err != nil {
		return err
	}

	ra, err := exec.RowsAffected()
	if err != nil {
		return err
	}

	if ra <= 0 {
		return ErrNoteNotFound
	}

	return nil
}
