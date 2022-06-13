package account

import (
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal/types"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

var registerLock sync.Mutex

func RegisterLocalAccount(
	email, password string,
) (*types.User, error) {
	registerLock.Lock()
	defer registerLock.Unlock()

	tx, err := database.Database.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
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

	_, err = tx.Exec("INSERT INTO user_email(id, email) VALUES ($1, $2)", id.String(), email)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO user_auth(id, password) VALUES ($1, $2)", id.String(), hashedPassword)
	if err != nil {
		return nil, err
	}

	var (
		hashedPasswordString = string(hashedPassword)
	)
	createdUser := types.User{
		ID:        id,
		Email:     &email,
		Password:  &hashedPasswordString,
		CreatedAt: createdAt,
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}
