package account

import (
	"context"
	"github.com/MatticNote/MatticNote/database"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"net/http"
)

func verifyMail(c *fiber.Ctx) error {
	token := c.Params("token")
	var targetUuid string
	var newEmail string

	err := database.DBPool.QueryRow(
		context.Background(),
		"select uuid, new_email from user_mail_transaction where token = $1 and expired_at > now();",
		token,
	).Scan(
		&targetUuid,
		&newEmail,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(http.StatusBadRequest).Render(
				"verify_error",
				fiber.Map{},
				"_layout/error",
			)
		} else {
			return err
		}
	}

	tx, err := database.DBPool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx) {
		_ = tx.Rollback(context.Background())
	}(tx)

	_, err = tx.Exec(
		context.Background(),
		"update user_mail set email = $1, is_verified = true where uuid = $2;",
		newEmail,
		targetUuid,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		"delete from user_mail_transaction where uuid = $1;",
		targetUuid,
	)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return c.Redirect("/account/login?verified=true")
}
