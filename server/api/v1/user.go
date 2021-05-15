package v1

import (
	"context"
	"github.com/MatticNote/MatticNote/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"net/http"
)

func getUser(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return v1BadRequest(c, "Not valid UUID format")
	}

	res := new(v1UserRes)
	var (
		isActive  bool
		isSuspend bool
	)

	err = database.DBPool.QueryRow(
		context.Background(),
		"select uuid, username, host, display_name, summary, created_at, updated_at, accept_manually, is_bot, is_active, is_suspend from \"user\" where uuid = $1;",
		targetUuid.String(),
	).Scan(
		&res.Uuid,
		&res.Username,
		&res.Host,
		&res.DisplayName,
		&res.Summary,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.AcceptManually,
		&res.IsBot,
		&isActive,
		&isSuspend,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return v1NotFound(c)
		} else {
			return err
		}
	}

	if !isActive {
		c.Status(http.StatusGone)
		return nil
	}

	if isSuspend {
		return v1Forbidden(c, "Specified user is suspended")
	}

	return c.Status(http.StatusOK).JSON(res)
}
