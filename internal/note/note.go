package note

import (
	"errors"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

//func CreateNote(
//	owner ksuid.KSUID,
//	cw,
//	body *string,
//) (*types.Note, error) {
//	newId := ksuid.New()
//	var (
//		createdAt sql.NullTime
//		cwNS      sql.NullString
//		bodyNS    sql.NullString
//	)
//	err := database.Database.QueryRow(
//		"INSERT INTO note(id, owner, cw, body) VALUES ($1, $2, $3, $4) RETURNING created_at, cw, body;",
//		newId.String(),
//		owner.String(),
//		cw,
//		body,
//	).Scan(
//		&createdAt,
//		&cwNS,
//		&bodyNS,
//	)
//
//	if err != nil {
//		return nil, err
//	}
//
//	newNote := &types.Note{
//		ID:        newId,
//		Owner:     owner,
//		CW:        cwNS,
//		Body:      bodyNS,
//		CreatedAt: createdAt,
//	}
//
//	return newNote, nil
//}
//
//func GetNote(id ksuid.KSUID) (*types.Note, error) {
//	note := new(types.Note)
//
//	err := database.Database.QueryRow(
//		"SELECT id, owner, cw, body, created_at FROM note WHERE id = $1",
//		id.String(),
//	).Scan(
//		&note.ID,
//		&note.Owner,
//		&note.CW,
//		&note.Body,
//		&note.CreatedAt,
//	)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, ErrNoteNotFound
//		} else {
//			return nil, err
//		}
//	}
//
//	return note, nil
//}
//
//func DeleteNote(id ksuid.KSUID) error {
//	exec, err := database.Database.Exec(
//		"DELETE FROM note WHERE id = $1",
//		id.String(),
//	)
//	if err != nil {
//		return err
//	}
//
//	ra, err := exec.RowsAffected()
//	if err != nil {
//		return err
//	}
//
//	if ra <= 0 {
//		return ErrNoteNotFound
//	}
//
//	return nil
//}
