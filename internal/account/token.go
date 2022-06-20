package account

import (
	"crypto/sha512"
	"database/sql"
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal/types"
	"github.com/segmentio/ksuid"
	"time"
)

var (
	ErrInvalidUserToken = errors.New("invalid user token")
)

func GenerateUserToken(
	userId ksuid.KSUID,
	issuedFromIP string,
	expiredAt ...time.Time,
) (string, error) {
	sessionId := sha512.Sum512([]byte(fmt.Sprintf("%s.%s", userId.String(), ksuid.New().String())))

	var expiredAtSQLNT sql.NullTime
	if len(expiredAt) > 0 {
		expiredAtSQLNT = sql.NullTime{
			Time:  expiredAt[0],
			Valid: true,
		}
	}

	sessionIdStr := fmt.Sprintf("%x", sessionId)

	_, err := database.Database.Exec(
		"INSERT INTO user_session(token, user_id, expired_at, issued_from) VALUES ($1, $2, $3, $4)",
		sessionIdStr,
		userId.String(),
		expiredAtSQLNT,
		issuedFromIP,
	)
	if err != nil {
		return "", err
	}

	return sessionIdStr, nil
}

func GetUserFromToken(
	token string,
) (*types.User, error) {
	var userId ksuid.KSUID
	if err := database.Database.QueryRow("SELECT user_id FROM user_session WHERE token = $1 AND (expired_at IS NULL OR expired_at >= now());", token).Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidUserToken
		} else {
			return nil, err
		}
	}

	user, err := GetUser(userId)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserGone
	}

	return user, nil
}

func DestroyUserToken(
	token string,
) error {
	_, err := database.Database.Exec("DELETE FROM user_session WHERE token = $1 OR expired_at < now();", token)
	if err != nil {
		return err
	}

	return nil
}
