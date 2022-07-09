package account

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/segmentio/ksuid"
	"sync"
	"time"
)

var (
	createInviteLock  sync.Mutex
	useInviteCodeLock sync.Mutex
)

var (
	ErrInvalidInviteCode = errors.New("invalid invite code")
)

const (
	inviteCodeLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	inviteCodeLength  = 8
)

func createInviteCode() (string, error) {
	b := make([]byte, inviteCodeLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	var result string
	for _, v := range b {
		result += string(inviteCodeLetters[int(v)%len(inviteCodeLetters)])
	}

	return result, nil
}

func CreateInvite(
	owner *ksuid.KSUID,
	count uint,
	expiredAt *time.Time,
) (*schemas.UserInvite, error) {
	createInviteLock.Lock()
	defer createInviteLock.Unlock()

	inviteInternalID := ksuid.New()
	inviteCode, err := createInviteCode()
	if err != nil {
		return nil, err
	}

	var countP *uint = nil
	if count > 0 {
		countP = &count
	}

	var invite = &schemas.UserInvite{
		ID:    inviteInternalID,
		Owner: owner,
		Code:  inviteCode,
	}

	err = database.Database.QueryRow(
		"INSERT INTO users_invite(id, owner, code, count, expired_at) VALUES ($1, $2, $3, $4, $5) RETURNING count, expired_at;",
		owner,
		inviteCode,
		countP,
		expiredAt,
	).Scan(
		&invite.Count,
		&invite.ExpiredAt,
	)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

func UseInviteCode(
	inviteCode string,
) error {
	useInviteCodeLock.Lock()
	defer useInviteCodeLock.Unlock()

	var (
		id    ksuid.KSUID
		count sql.NullInt32
	)

	err := database.Database.QueryRow(
		"SELECT id, count FROM users_invite WHERE code = $1 AND (expired_at IS NULL OR expired_at >= now()) AND count > 0",
		inviteCode,
	).Scan(
		&id,
		&count,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidInviteCode
		} else {
			return err
		}
	}

	if count.Valid {
		newCountRaw, err := count.Value()
		if err != nil {
			return err
		}
		i, ok := newCountRaw.(int32)
		if !ok {
			return errors.New("internal error")
		}
		newCount := i - 1

		_, err = database.Database.Exec(
			"UPDATE users_invite SET count=$1 WHERE id = $2",
			newCount,
			id,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
