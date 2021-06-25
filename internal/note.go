package internal

import (
	"context"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

var (
	ErrNoteNotFound = errors.New("specified note is not found")
)

type NoteStruct struct {
	Uuid      uuid.UUID
	Author    UserStruct
	CreatedAt misc.NullTime
	Cw        misc.NullString
	Body      misc.NullString
	LocalOnly bool
}

func CreateNoteFromLocal(authorUuid uuid.UUID, cw, text string, replyUuid, reTextUuid uuid.UUID, localOnly bool) (*NoteStruct, error) {
	newNoteUuid := uuid.Must(uuid.NewUUID())

	tx, err := database.DBPool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer func(tx pgx.Tx) {
		_ = tx.Rollback(context.Background())
	}(tx)

	_, err = tx.Exec(
		context.Background(),
		"insert into note(uuid, author, local_only) VALUES ($1, $2, $3);",
		newNoteUuid,
		authorUuid,
		localOnly,
	)
	if err != nil {
		return nil, err
	}

	if cw != "" {
		_, err = tx.Exec(
			context.Background(),
			"update note set cw = $1 where uuid = $2;",
			cw,
			newNoteUuid,
		)
		if err != nil {
			return nil, err
		}
	}

	if text != "" {
		_, err = tx.Exec(
			context.Background(),
			"update note set body = $1 where uuid = $2;",
			text,
			newNoteUuid,
		)
		if err != nil {
			return nil, err
		}
	}

	if replyUuid != uuid.Nil {
		_, err = tx.Exec(
			context.Background(),
			"update note set reply_uuid = $1 where uuid = $2;",
			replyUuid,
			newNoteUuid,
		)
		if err != nil {
			return nil, ErrNoteNotFound
		}
	}

	if reTextUuid != uuid.Nil {
		_, err = tx.Exec(
			context.Background(),
			"update note set retext_uuid = $1 where uuid = $2;",
			reTextUuid,
			newNoteUuid,
		)
		if err != nil {
			return nil, ErrNoteNotFound
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	// TODO: ActivityPubのCreateアクティビティを書く

	createdNote, err := GetNote(newNoteUuid)
	if err != nil {
		return nil, err
	}

	return createdNote, err
}

func GetNote(noteUuid uuid.UUID) (*NoteStruct, error) {
	var (
		noteRes       NoteStruct
		noteAuthorRes UserStruct
		isActive      bool
		isSuspend     bool
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"select note.uuid, note.created_at, cw, body, local_only,"+
			" u.uuid, username, host, display_name, summary, u.created_at, updated_at, accept_manually, is_bot, is_active, is_suspend "+
			"from note, \"user\" u where u.uuid = note.author and note.uuid = $1;",
		noteUuid.String(),
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
			return nil, ErrNoteNotFound
		} else {
			return nil, err
		}
	}

	if !isActive {
		return nil, ErrUserGone
	}

	if isSuspend {
		return nil, ErrUserSuspended
	}

	noteRes.Author = noteAuthorRes
	return &noteRes, nil
}

func DeleteNote(noteUuid uuid.UUID) error {
	exec, err := database.DBPool.Exec(
		context.Background(),
		"delete from note where uuid = $1;",
		noteUuid.String(),
	)
	if err != nil {
		return err
	}

	if exec.RowsAffected() <= 0 {
		return ErrNoteNotFound
	}

	// TODO: ActivityPubに関わるDeleteアクティビティなどを書く

	return nil
}
