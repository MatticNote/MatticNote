package account

import (
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/lib/pq"
	"github.com/segmentio/ksuid"
)

var (
	ErrCannotRelationYourself = errors.New("you can't change your own relationship")
	ErrAlreadyFollowing       = errors.New("you're already following")
	ErrNotFollowing           = errors.New("you are not following")
)

func CreateFollowRelation(
	fromUser ksuid.KSUID,
	toUser ksuid.KSUID,
) (bool, error) {
	if fromUser == toUser {
		return false, ErrCannotRelationYourself
	}

	_, err := database.Database.Exec(
		"INSERT INTO users_follow_relation(from_follow, to_follow, is_active) VALUES ($1, $2, $3)",
		fromUser.String(),
		toUser.String(),
		true, // TODO: Follow lock system
	)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Code.Name() {
			case "check_violation":
				return false, ErrCannotRelationYourself
			case "unique_violation":
				return false, ErrAlreadyFollowing
			case "foreign_key_violation":
				return false, ErrUserNotFound
			default:
				return false, err
			}
		} else {
			return false, err
		}
	}

	return true, err
}

func DeleteFollowRelation(
	fromUser ksuid.KSUID,
	toUser ksuid.KSUID,
) error {
	if fromUser == toUser {
		return ErrCannotRelationYourself
	}

	exec, err := database.Database.Exec(
		"DELETE FROM users_follow_relation WHERE from_follow=$1 AND to_follow=$2",
		fromUser.String(),
		toUser.String(),
	)
	if err != nil {
		return err
	}

	ra, err := exec.RowsAffected()
	if err != nil {
		return err
	}

	if ra <= 0 {
		return ErrNotFollowing
	}

	return nil
}

func ListFollowingRelation(
	userId ksuid.KSUID,
	limit int,
	offset int,
) ([]*schemas.User, error) {
	rows, err := database.Database.Query(
		"SELECT to_follow FROM users_follow_relation WHERE from_follow = $1 AND is_active IS TRUE ORDER BY following_since DESC LIMIT $2 OFFSET $3;",
		userId,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var followingIds []ksuid.KSUID
	for rows.Next() {
		var id ksuid.KSUID
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		followingIds = append(followingIds, id)
	}

	return GetUserMultiple(followingIds...)
}

func ListFollowerRelation(
	userId ksuid.KSUID,
	limit int,
	offset int,
) ([]*schemas.User, error) {
	rows, err := database.Database.Query(
		"SELECT from_follow FROM users_follow_relation WHERE to_follow = $1 AND is_active IS TRUE ORDER BY following_since DESC LIMIT $2 OFFSET $3;",
		userId,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var followingIds []ksuid.KSUID
	for rows.Next() {
		var id ksuid.KSUID
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		followingIds = append(followingIds, id)
	}

	return GetUserMultiple(followingIds...)
}
