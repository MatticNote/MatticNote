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

type (
	RelationStruct struct {
		ID     ksuid.KSUID
		userId ksuid.KSUID
		User   *schemas.User
	}
)

func CreateFollowRelation(
	fromUser ksuid.KSUID,
	toUser ksuid.KSUID,
) (bool, error) {
	if fromUser == toUser {
		return false, ErrCannotRelationYourself
	}

	relationId := ksuid.New()

	_, err := database.Database.Exec(
		"INSERT INTO users_follow_relation(id, from_follow, to_follow, is_active) VALUES ($1, $2, $3, $4)",
		relationId.String(),
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
	maxId ksuid.KSUID,
	sinceId ksuid.KSUID,
	limit int,
) ([]*RelationStruct, error) {
	rows, err := database.Database.Query(
		"SELECT id, to_follow FROM users_follow_relation WHERE from_follow = $1 AND is_active IS TRUE AND id >= $2 AND id < $3 ORDER BY id DESC LIMIT $4",
		userId,
		sinceId.String(),
		maxId.String(),
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var followingIds []ksuid.KSUID
	var relationInfo []*RelationStruct
	for rows.Next() {
		relation := new(RelationStruct)

		var (
			followUserId ksuid.KSUID
		)
		err := rows.Scan(&relation.ID, &followUserId)
		if err != nil {
			return nil, err
		}
		relation.userId = followUserId

		followingIds = append(followingIds, followUserId)
		relationInfo = append(relationInfo, relation)
	}

	followingUsers, err := GetUserMultiple(followingIds...)
	if err != nil {
		return nil, err
	}

	followingUserDict := make(map[string]*schemas.User)
	for _, v := range followingUsers {
		followingUserDict[v.ID.String()] = v
	}

	for i, v := range relationInfo {
		user, exists := followingUserDict[v.userId.String()]
		if !exists {
			continue
		}
		relationInfo[i].User = user
	}

	return relationInfo, nil
}

func ListFollowerRelation(
	userId ksuid.KSUID,
	maxId ksuid.KSUID,
	sinceId ksuid.KSUID,
	limit int,
) ([]*RelationStruct, error) {
	rows, err := database.Database.Query(
		"SELECT id, from_follow FROM users_follow_relation WHERE to_follow = $1 AND is_active IS TRUE AND id >= $2 AND id < $3 ORDER BY id DESC LIMIT $4",
		userId,
		sinceId.String(),
		maxId.String(),
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var followerIds []ksuid.KSUID
	var relationInfo []*RelationStruct
	for rows.Next() {
		relation := new(RelationStruct)

		var (
			followUserId ksuid.KSUID
		)
		err := rows.Scan(&relation.ID, &followUserId)
		if err != nil {
			return nil, err
		}
		relation.userId = followUserId

		followerIds = append(followerIds, followUserId)
		relationInfo = append(relationInfo, relation)
	}

	followingUsers, err := GetUserMultiple(followerIds...)
	if err != nil {
		return nil, err
	}

	followingUserDict := make(map[string]*schemas.User)
	for _, v := range followingUsers {
		followingUserDict[v.ID.String()] = v
	}

	for i, v := range relationInfo {
		user, exists := followingUserDict[v.userId.String()]
		if !exists {
			continue
		}
		relationInfo[i].User = user
	}

	return relationInfo, nil
}
