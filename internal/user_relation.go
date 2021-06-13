package internal

import (
	"context"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/google/uuid"
)

var (
	ErrAlreadyFollowing   = errors.New("the specified user is already following")
	ErrNotFollowing       = errors.New("the specified user is not following")
	ErrAlreadyBlocking    = errors.New("the specified user is already blocking")
	ErrNotBlocking        = errors.New("the specified user is not blocking")
	ErrTargetBlocked      = errors.New("the specified user is blocking")
	ErrCantRelateYourself = errors.New("can't yourself")
)

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
