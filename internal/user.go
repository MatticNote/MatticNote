package internal

import (
	"context"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	ErrUserExists        = errors.New("username or email is already taken")
	ErrLoginFailed       = errors.New("there is an error in the login name or password")
	ErrEmailAuthRequired = errors.New("target user required email authentication")
	Err2faRequired       = errors.New("target user required two factor authentication") // WIP
)

const (
	UserPasswordHashCost = 12
	UserKeyPairLength    = 2048
)

func RegisterLocalUser(email, username, password string, skipEmailVerify bool) (uuid.UUID, error) {
	var count int
	err := database.DBPool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM \"user\" LEFT JOIN user_mail um on \"user\".uuid = um.uuid WHERE (username ILIKE $1 AND host IS NULL) OR email ILIKE $2;",
		username,
		email,
	).Scan(&count)
	if err != nil && err != pgx.ErrNoRows {
		return uuid.Nil, err
	}

	if count > 0 {
		return uuid.Nil, ErrUserExists
	}

	newUuid := uuid.Must(uuid.NewRandom())

	tx, err := database.DBPool.Begin(context.Background())
	if err != nil {
		return uuid.Nil, err
	}
	defer func(tx pgx.Tx) {
		_ = tx.Rollback(context.Background())
	}(tx)

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO \"user\"(uuid, username) VALUES ($1, $2);",
		newUuid,
		username,
	)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO user_mail(uuid, email, is_verified) VALUES ($1, $2, $3);",
		newUuid,
		email,
		skipEmailVerify,
	)
	if err != nil {
		return uuid.Nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), UserPasswordHashCost)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO user_password(uuid, password) VALUES ($1, $2);",
		newUuid,
		hashedPassword,
	)
	if err != nil {
		return uuid.Nil, err
	}

	rsaPrivateKey, rsaPublicKey := misc.GenerateRSAKeypair(UserKeyPairLength)

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO user_signature_key(uuid, public_key, private_key) VALUES ($1, $2, $3);",
		newUuid,
		string(rsaPublicKey),
		string(rsaPrivateKey),
	)
	if err != nil {
		return uuid.Nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return uuid.Nil, err
	}

	if !skipEmailVerify {
		// TODO: メールを認証するための処理
	}

	return newUuid, err
}

func ValidateLoginUser(login, password string) (uuid.UUID, error) {
	var (
		targetUuid     uuid.UUID
		isMailVerified bool
		targetPassword []byte
	)

	err := database.DBPool.QueryRow(
		context.Background(),
		"SELECT \"user\".uuid, is_verified, password FROM \"user\", user_mail um, user_password up WHERE \"user\".uuid = um.uuid AND \"user\".uuid = up.uuid AND email ILIKE $1 OR (username ILIKE $1 AND host IS NULL);",
		login,
	).Scan(
		&targetUuid,
		&isMailVerified,
		&targetPassword,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, ErrLoginFailed
		} else {
			return uuid.Nil, err
		}
	}

	if err := bcrypt.CompareHashAndPassword(targetPassword, []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return uuid.Nil, ErrLoginFailed
		} else {
			return uuid.Nil, err
		}
	}

	if !isMailVerified {
		return uuid.Nil, ErrEmailAuthRequired
	}

	return targetUuid, nil
}
