package internal

import (
	"context"
	"github.com/MatticNote/MatticNote/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func CreateNoteFromLocal(authorUuid uuid.UUID, cw, text string, replyUuid, reTextUuid uuid.UUID, localOnly bool) (uuid.UUID, error) {
	newNoteUuid := uuid.Must(uuid.NewUUID())

	tx, err := database.DBPool.Begin(context.Background())
	if err != nil {
		return uuid.Nil, err
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
		return uuid.Nil, err
	}

	if cw != "" {
		_, err = tx.Exec(
			context.Background(),
			"update note set cw = $1 where uuid = $2;",
			cw,
			newNoteUuid,
		)
		if err != nil {
			return uuid.Nil, err
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
			return uuid.Nil, err
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
			return uuid.Nil, err
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
			return uuid.Nil, err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return uuid.Nil, err
	}
	return newNoteUuid, err
}
