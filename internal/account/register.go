package account

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal/types"
	"github.com/gomodule/redigo/redis"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

var (
	registerLock       sync.Mutex
	chooseUsernameLock sync.Mutex
)

var (
	ErrInvalidToken          = errors.New("token is not valid")
	ErrUsernameAlreadyTaken  = errors.New("username is already taken")
	ErrUsernameAlreadyChosen = errors.New("username is already chosen")
)

func RegisterLocalAccount(
	email, password string,
	skipEmailVerify bool,
) (*types.User, error) {
	registerLock.Lock()
	defer registerLock.Unlock()

	tx, err := database.Database.Begin()
	if err != nil {
		return nil, err
	}
	rsCon := database.RedisPool.Get()
	defer func() {
		_ = tx.Rollback()
		_ = rsCon.Close()
	}()

	var createdAt time.Time
	id := ksuid.New()
	err = tx.QueryRow("INSERT INTO \"user\"(id) VALUES ($1) RETURNING created_at", id.String()).Scan(&createdAt)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO user_email(id, email, verified) VALUES ($1, $2, $3)", id.String(), email, skipEmailVerify)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO user_auth(id, password) VALUES ($1, $2)", id.String(), string(hashedPassword))
	if err != nil {
		return nil, err
	}

	if !skipEmailVerify {
		verifyKey := fmt.Sprintf("emailVerify:%s", ksuid.New().String())
		_, err = rsCon.Do(
			"HSET",
			verifyKey,
			"id",
			id.String(),
			"email",
			email,
		)
		if err != nil {
			return nil, err
		}
		err := rsCon.Send("EXPIRE", verifyKey, "3600")
		if err != nil {
			return nil, err
		}
		// TODO: Send verification mail
	}

	createdUser := types.User{
		ID: id,
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		CreatedAt: sql.NullTime{
			Time:  createdAt,
			Valid: true,
		},
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func VerifyEmailToken(
	token string,
) error {
	rsCon := database.RedisPool.Get()
	defer func() {
		_ = rsCon.Close()
	}()

	tokenKey := fmt.Sprintf("emailVerify:%s", token)
	userIdStr, err := redis.String(rsCon.Do("HGET", tokenKey, "id"))
	if err != nil {
		if err == redis.ErrNil {
			return ErrInvalidToken
		} else {
			return err
		}
	}

	userId, err := ksuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	email, err := redis.String(rsCon.Do("HGET", tokenKey, "email"))
	if err != nil {
		return err
	}

	_, err = database.Database.Exec("UPDATE user_email SET verified=TRUE, email=$1 WHERE id = $2;", email, userId.String())
	if err != nil {
		return err
	}

	err = rsCon.Send("DEL", tokenKey)
	if err != nil {
		return err
	}

	return nil
}

func ChooseUsername(
	userId ksuid.KSUID,
	username string,
) error {
	chooseUsernameLock.Lock()
	defer chooseUsernameLock.Unlock()

	var exists int
	err := database.Database.QueryRow("SELECT count(*) AS exists FROM \"user\" WHERE username = $1 AND host IS NULL", username).Scan(&exists)
	if err != nil {
		return err
	}

	if exists > 0 {
		return ErrUsernameAlreadyTaken
	}

	exec, err := database.Database.Exec("UPDATE \"user\" SET username=$1 WHERE username IS NULL AND host IS NULL AND id=$2", username, userId.String())
	if err != nil {
		return err
	}

	ra, err := exec.RowsAffected()
	if err != nil {
		return err
	}

	if ra <= 0 {
		return ErrUsernameAlreadyChosen
	}

	return nil
}
