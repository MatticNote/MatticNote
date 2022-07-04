package account

import (
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUser(user *schemas.User) error {
	_, err := database.Database.Exec(
		"UPDATE users SET display_name=$1, headline=$2, description=$3 WHERE id = $4",
		user.DisplayName,
		user.Headline,
		user.Description,
		user.ID.String(),
	)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserPassword(userId ksuid.KSUID, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	exec, err := database.Database.Exec(
		"UPDATE users_auth SET password = $1 WHERE id = $2",
		hashedPassword,
		userId.String(),
	)
	if err != nil {
		return err
	}

	ra, err := exec.RowsAffected()
	if err != nil {
		return err
	}

	if ra <= 0 {
		return ErrUserNotFound
	}

	return nil
}
