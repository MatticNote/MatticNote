package account

import (
	"crypto/sha512"
	"database/sql"
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/segmentio/ksuid"
	"time"
)

var (
	ErrInvalidUserToken  = errors.New("invalid user token")
	ErrUserTokenNotFound = errors.New("user token not found")
)

func GenerateUserToken(
	userId ksuid.KSUID,
	issuedFromIP string,
	expiredAt ...time.Time,
) (string, error) {
	tokenRaw := sha512.Sum512([]byte(fmt.Sprintf("%s.%s", userId.String(), ksuid.New())))
	token := fmt.Sprintf("%x", tokenRaw)

	var err error

	if len(expiredAt) > 0 {
		_, err = database.Database.Exec(
			"INSERT INTO users_token(id, token, user_id, expired_at, ip) VALUES ($1, $2, $3, $4, $5)",
			ksuid.New().String(),
			token,
			userId.String(),
			expiredAt[0],
			issuedFromIP,
		)
	} else {
		_, err = database.Database.Exec(
			"INSERT INTO users_token(id, token, user_id, ip) VALUES ($1, $2, $3, $4)",
			ksuid.New().String(),
			token,
			userId.String(),
			issuedFromIP,
		)
	}

	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUserFromToken(
	token string,
) (*schemas.User, error) {
	var userId ksuid.KSUID

	err := database.Database.QueryRow("SELECT user_id FROM users_token WHERE token = $1", token).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidUserToken
		} else {
			return nil, err
		}
	}

	user, err := GetUser(userId)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt.Valid && user.DeletedAt.Time.Before(time.Now()) {
		return nil, ErrUserGone
	}

	return user, nil
}

func GetUserToken(
	id ksuid.KSUID,
) (*schemas.UserToken, error) {
	token := new(schemas.UserToken)

	err := database.Database.QueryRow(
		"SELECT id, token, user_id, expired_at, ip, created_at FROM users_token WHERE id = $1 AND (expired_at IS NULL OR expired_at >= now());",
		id.String(),
	).Scan(
		&token.ID,
		&token.Token,
		&token.UserId,
		&token.ExpiredAt,
		&token.IP,
		&token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserTokenNotFound
		} else {
			return nil, err
		}
	}

	return token, nil
}

func ListUserToken(
	userId ksuid.KSUID,
) ([]*schemas.UserToken, error) {
	var tokenList []*schemas.UserToken

	row, err := database.Database.Query("SELECT id, token, user_id, expired_at, ip, created_at FROM users_token WHERE user_id = $1 AND (expired_at IS NULL OR expired_at >= now());", userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tokenList, nil
		} else {
			return nil, err
		}
	}
	defer func() {
		_ = row.Close()
	}()

	for row.Next() {
		token := new(schemas.UserToken)
		if err := row.Scan(&token.ID, &token.Token, &token.UserId, &token.ExpiredAt, &token.IP, &token.CreatedAt); err != nil {
			return nil, err
		}
		tokenList = append(tokenList, token)
	}

	return tokenList, err
}

func DeleteUserToken(
	id ksuid.KSUID,
) error {
	_, err := database.Database.Exec("DELETE FROM users_token WHERE id = $1 OR expired_at < now();", id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserTokenFromToken(
	token string,
) error {
	_, err := database.Database.Exec("DELETE FROM users_token WHERE token = $1 OR expired_at < now();", token)
	if err != nil {
		return err
	}

	return nil
}
