package internal

import (
	"context"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

var (
	ErrAlreadyFollowing   = errors.New("the specified user is already following")
	ErrNotFollowing       = errors.New("the specified user is not following")
	ErrAlreadyMuting      = errors.New("the specified user is already muting")
	ErrNotMuting          = errors.New("the specified user is not muting")
	ErrAlreadyBlocking    = errors.New("the specified user is already blocking")
	ErrNotBlocking        = errors.New("the specified user is not blocking")
	ErrTargetBlocked      = errors.New("the specified user is blocking")
	ErrCantRelateYourself = errors.New("can't yourself")
	ErrUnknownRequest     = errors.New("unknown follow request")
)

var emptyUuidList = make([]uuid.UUID, 0)

func CreateFollowRelation(fromUser, targetUser uuid.UUID, isPending bool) error {
	if fromUser == targetUser {
		return ErrCantRelateYourself
	}

	var isBlocked int

	err := database.DBPool.QueryRow(
		context.Background(),
		"select count(*) from block_relation where (block_from = $2 and block_to = $1) or (block_from = $1 and block_to = $2);",
		fromUser.String(),
		targetUser.String(),
	).Scan(&isBlocked)
	if err != nil {
		return err
	}
	if isBlocked > 0 {
		return ErrTargetBlocked
	}

	exec, err := database.DBPool.Exec(
		context.Background(),
		"insert into follow_relation(follow_from, follow_to, is_pending) values ($1, $2, $3) on conflict do nothing;",
		fromUser.String(),
		targetUser.String(),
		isPending,
	)
	if err != nil {
		return err
	}
	if exec.RowsAffected() <= 0 {
		return ErrAlreadyFollowing
	}

	return nil
}

func DestroyFollowRelation(fromUser, targetUser uuid.UUID) error {
	if fromUser == targetUser {
		return ErrCantRelateYourself
	}

	exec, err := database.DBPool.Exec(
		context.Background(),
		"delete from follow_relation where follow_from = $1 and follow_to = $2;",
		fromUser.String(),
		targetUser.String(),
	)
	if err != nil {
		return err
	}
	if exec.RowsAffected() <= 0 {
		return ErrNotFollowing
	}

	return nil
}

func CreateMuteRelation(fromUser, targetUser uuid.UUID) error {
	if fromUser == targetUser {
		return ErrCantRelateYourself
	}

	exec, err := database.DBPool.Exec(
		context.Background(),
		"insert into mute_relation(mute_from, mute_to) values ($1, $2) on conflict do nothing;",
		fromUser.String(),
		targetUser.String(),
	)
	if err != nil {
		return err
	}
	if exec.RowsAffected() <= 0 {
		return ErrAlreadyMuting
	}

	return nil
}

func DestroyMuteRelation(fromUser, targetUser uuid.UUID) error {
	if fromUser == targetUser {
		return ErrCantRelateYourself
	}

	exec, err := database.DBPool.Exec(
		context.Background(),
		"delete from mute_relation where mute_from = $1 and mute_to = $2;",
		fromUser.String(),
		targetUser.String(),
	)
	if err != nil {
		return err
	}
	if exec.RowsAffected() <= 0 {
		return ErrNotMuting
	}

	return nil
}

func CreateBlockRelation(fromUser, targetUser uuid.UUID) error {
	if fromUser == targetUser {
		return ErrCantRelateYourself
	}

	tx, err := database.DBPool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	_, err = tx.Exec(
		context.Background(),
		"delete from follow_relation where (follow_from = $1 and follow_to = $2) or (follow_from = $2 and follow_to = $1);",
		fromUser.String(),
		targetUser.String(),
	)
	if err != nil {
		return err
	}

	exec, err := tx.Exec(
		context.Background(),
		"insert into block_relation(block_from, block_to) values ($1, $2) on conflict do nothing;",
		fromUser.String(),
		targetUser.String(),
	)
	if err != nil {
		return err
	}
	if exec.RowsAffected() <= 0 {
		return ErrAlreadyBlocking
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}

func DestroyBlockRelation(fromUser, targetUser uuid.UUID) error {
	if fromUser == targetUser {
		return ErrCantRelateYourself
	}

	exec, err := database.DBPool.Exec(
		context.Background(),
		"delete from block_relation where block_from = $1 and block_to = $2;",
		fromUser.String(),
		targetUser.String(),
	)
	if err != nil {
		return err
	}
	if exec.RowsAffected() <= 0 {
		return ErrNotBlocking
	}

	return nil
}

func ListFollowRequests(fromUser uuid.UUID) ([]uuid.UUID, error) {
	rows, err := database.DBPool.Query(
		context.Background(),
		"select follow_from from follow_relation where follow_to = $1 and is_pending = true;",
		fromUser.String(),
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return emptyUuidList, nil
		} else {
			return emptyUuidList, err
		}
	}
	defer rows.Close()

	var requests = make([]uuid.UUID, 0)
	for rows.Next() {
		var reqUuid uuid.UUID
		err = rows.Scan(&reqUuid)
		if err != nil {
			return emptyUuidList, err
		}
		requests = append(requests, reqUuid)
	}
	if rows.Err() != nil {
		return emptyUuidList, rows.Err()
	}

	return requests, nil
}

func AcceptFollowRequest(fromUser, targetUser uuid.UUID) error {
	if fromUser == targetUser {
		return ErrCantRelateYourself
	}

	exec, err := database.DBPool.Exec(
		context.Background(),
		"update follow_relation set is_pending = false where follow_to = $1 and follow_from = $2 and is_pending = true;",
		fromUser.String(),
		targetUser.String(),
	)
	if err != nil {
		return err
	}
	if exec.RowsAffected() <= 0 {
		return ErrUnknownRequest
	}

	// TODO: 通知とか飛ばせるようにしたい

	return nil
}

func RejectFollowRequest(fromUser, targetUser uuid.UUID) error {
	if fromUser == targetUser {
		return ErrCantRelateYourself
	}

	exec, err := database.DBPool.Exec(
		context.Background(),
		"delete from follow_relation where follow_to = $1 and follow_from = $2 and is_pending = true;",
		fromUser.String(),
		targetUser.String(),
	)
	if err != nil {
		return err
	}
	if exec.RowsAffected() <= 0 {
		return ErrUnknownRequest
	}

	return nil
}
