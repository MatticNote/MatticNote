package internal

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists = errors.New("username or email is already taken")
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

	rsaKeyRaw, err := rsa.GenerateKey(rand.Reader, UserKeyPairLength)
	if err != nil {
		return uuid.Nil, err
	}
	rsaPrivateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rsaKeyRaw),
	})
	rsaPublicKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(rsaKeyRaw.Public().(*rsa.PublicKey)),
	})

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
