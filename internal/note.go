package internal

import (
	"context"
	"errors"
	"github.com/MatticNote/MatticNote/activitypub"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/MatticNote/MatticNote/worker"
	"github.com/gocraft/work"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"strings"
)

var (
	ErrNoteNotFound      = errors.New("specified note is not found")
	ErrInvalidVisibility = errors.New("invalid note visibility")
)

type NoteStruct struct {
	Uuid       uuid.UUID
	Author     *UserStruct
	CreatedAt  misc.NullTime
	Cw         misc.NullString
	Body       misc.NullString
	LocalOnly  bool
	Reply      *NoteStruct
	ReText     *NoteStruct
	Visibility string
}

func CreateNoteFromLocal(authorUuid uuid.UUID, cw, text *string, replyUuid, reTextUuid *uuid.UUID, localOnly bool, visibility string) (*NoteStruct, error) {
	switch strings.ToUpper(visibility) {
	case "PUBLIC":
	case "UNLISTED":
	case "FOLLOWER":
	case "DIRECT":
	default:
		return nil, ErrInvalidVisibility
	}

	newNoteUuid := uuid.Must(uuid.NewUUID())

	tx, err := database.DBPool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	_, err = tx.Exec(
		context.Background(),
		"insert into note(uuid, author, local_only, visibility) VALUES ($1, $2, $3, $4);",
		newNoteUuid,
		authorUuid,
		localOnly,
		strings.ToUpper(visibility),
	)
	if err != nil {
		return nil, err
	}

	if cw != nil {
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

	if text != nil {
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

	if replyUuid != nil {
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

	if reTextUuid != nil {
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

	createdNote, err := GetNote(newNoteUuid)
	if err != nil {
		return nil, err
	}

	if !localOnly {
		activity := activitypub.RenderNoteActivity(createdNote)
		switch strings.ToUpper(visibility) {
		case "PUBLIC", "UNLISTED", "FOLLOWER":
			followerInbox, err := GetUserFollowerInbox(authorUuid)
			if err != nil {
				return nil, err
			}
			if len(followerInbox) > 0 {
				for _, inbox := range followerInbox {
					_, err = worker.Enqueue.Enqueue(
						worker.JobDeliver,
						work.Q{
							"to":       inbox,
							"body":     activity,
							"fromUuid": authorUuid,
						},
					)
				}

			}
		case "DIRECT":
			// todo: ダイレクトによるdeliver
		}
	}

	return createdNote, err
}

func GetNote(noteUuid uuid.UUID, recursive ...bool) (*NoteStruct, error) {
	var (
		noteRes    NoteStruct
		authorUuid uuid.UUID
		replyUuid  uuid.UUID
		reTextUuid uuid.UUID
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"select uuid, created_at, cw, body, local_only, visibility, author, reply_uuid, retext_uuid "+
			"from note where note.uuid = $1;",
		noteUuid.String(),
	).Scan(
		&noteRes.Uuid,
		&noteRes.CreatedAt,
		&noteRes.Cw,
		&noteRes.Body,
		&noteRes.LocalOnly,
		&noteRes.Visibility,
		&authorUuid,
		&replyUuid,
		&reTextUuid,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoteNotFound
		} else {
			return nil, err
		}
	}

	noteRes.Author, err = GetUser(authorUuid)
	if err != nil {
		return nil, err
	}

	if noteRes.Author.IsSuspend {
		return nil, ErrUserSuspended
	}

	if len(recursive) > 0 && recursive[0] {
		if replyUuid != uuid.Nil {
			noteRes.Reply, err = GetNote(replyUuid, false)
			if err != nil {
				return nil, err
			}
		}
		if reTextUuid != uuid.Nil {
			noteRes.ReText, err = GetNote(reTextUuid, false)
			if err != nil {
				return nil, err
			}
		}
	}

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
