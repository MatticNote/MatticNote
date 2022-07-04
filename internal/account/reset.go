package account

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/gomodule/redigo/redis"
	"github.com/segmentio/ksuid"
	"sync"
)

var (
	newPasswordResetTokenLock     sync.Mutex
	isValidPasswordResetTokenLock sync.Mutex
	usePasswordResetTokenLock     sync.Mutex
)

var (
	ErrInvalidPasswordResetToken = errors.New("invalid password token")
)

func NewPasswordResetToken(
	email string,
) error {
	newPasswordResetTokenLock.Lock()
	defer newPasswordResetTokenLock.Unlock()

	rsCon := database.RedisPool.Get()
	defer func() {
		_ = rsCon.Close()
	}()

	var (
		userId    ksuid.KSUID
		userEmail string
	)

	err := database.Database.QueryRow(
		"SELECT u.id, email FROM users_email LEFT OUTER JOIN users u on u.id = users_email.id WHERE email ILIKE $1 AND is_verified IS TRUE AND (u.deleted_at IS NULL OR u.deleted_at >= now())",
		email,
	).Scan(&userId, &userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	resetKey := ksuid.New()
	resetRedisKey := fmt.Sprintf("passwordReset:%s", resetKey.String())

	_, err = rsCon.Do(
		"SET",
		resetRedisKey,
		userId.String(),
	)
	if err != nil {
		return err
	}
	err = rsCon.Send("EXPIRE", resetRedisKey, "900")
	if err != nil {
		return err
	}
	// TODO: Send verification mail

	return nil
}

func IsValidPasswordResetToken(token string) bool {
	isValidPasswordResetTokenLock.Lock()
	defer isValidPasswordResetTokenLock.Unlock()

	rsCon := database.RedisPool.Get()
	defer func() {
		_ = rsCon.Close()
	}()

	tokenRedisKey := fmt.Sprintf("passwordReset:%s", token)
	_, err := redis.String(rsCon.Do("GET", tokenRedisKey))
	if err != nil {
		if err == redis.ErrNil {
			return false
		} else {
			return false
		}
	}

	return true
}

func PopPasswordResetToken(token string) (*ksuid.KSUID, error) {
	usePasswordResetTokenLock.Lock()
	defer usePasswordResetTokenLock.Unlock()

	rsCon := database.RedisPool.Get()
	defer func() {
		_ = rsCon.Close()
	}()

	tokenRedisKey := fmt.Sprintf("passwordReset:%s", token)

	userIdStr, err := redis.String(rsCon.Do("GET", tokenRedisKey))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return nil, ErrInvalidPasswordResetToken
		} else {
			return nil, err
		}
	}

	userId, err := ksuid.Parse(userIdStr)
	if err != nil {
		return nil, err
	}

	ra, err := redis.Int(rsCon.Do("DEL", tokenRedisKey))
	if err != nil {
		return nil, err
	}

	if ra <= 0 {
		return nil, ErrInvalidPasswordResetToken
	}

	return &userId, nil
}
