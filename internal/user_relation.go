package internal

import (
	"context"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/google/uuid"
)

var (
	ErrAlreadyFollowing   = errors.New("the specified user is already following")
	ErrTargetBlocked      = errors.New("the specified user is blocked")
	ErrCantFollowYourself = errors.New("can't follow yourself")
)

func CreateFollowRelation(fromUser, targetUser uuid.UUID) error {
	if fromUser == targetUser {
		return ErrCantFollowYourself
	}

	var isBlocked int

	err := database.DBPool.QueryRow(
		context.Background(),
		"select count(*) from block_relation where block_from = $2 and block_to = $1;",
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
		"insert into follow_relation(follow_from, follow_to) values ($1, $2) on conflict do nothing;",
		fromUser.String(),
		targetUser.String(),
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
	// TODO: つくる
	return nil
}
