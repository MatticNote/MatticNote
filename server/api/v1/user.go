package v1

import (
	"context"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
)

func getUser(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return badRequest(c, "Not valid UUID format")
	}

	currentUsr, ok := c.Locals(loginUserLocal).(*internal.LocalUserStruct)
	if ok {
		log.Println(currentUsr.Uuid.String())
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
			return notFound(c, "Specified user was not found")
		} else {
			return err
		}
	}

	if !isActive {
		c.Status(http.StatusGone)
		return nil
	}

	if isSuspend {
		return forbidden(c, "Specified user is suspended")
	}

	return c.Status(http.StatusOK).JSON(res)
}
