package account

import (
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/database/schemas"
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
