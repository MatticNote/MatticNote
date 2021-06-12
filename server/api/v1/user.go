package v1

import (
	"context"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func getUser(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return badRequest(c, "Not valid UUID format")
	}

	// TODO: ターゲットのアカウントリレーション次第では403とかを出す
	//currentUsr, ok := c.Locals(internal.LoginUserLocal).(*internal.LocalUserStruct)

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
		c.Status(fiber.StatusGone)
		return nil
	}

	if isSuspend {
		return forbidden(c, "Specified user is suspended")
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func followUser(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return badRequest(c, "Not valid UUID format")
	}

	currentUsr, ok := c.Locals(internal.LoginUserLocal).(*internal.LocalUserStruct)
	if !ok {
		return unauthorized(c)
	}

	var (
		isActive  bool
		isSuspend bool
	)

	err = database.DBPool.QueryRow(
		context.Background(),
		"select is_active, is_suspend from \"user\" where uuid = $1;",
		targetUuid.String(),
	).Scan(
		&isActive,
		&isSuspend,
	)

	if !isActive {
		c.Status(fiber.StatusGone)
		return nil
	}

	if isSuspend {
		c.Status(fiber.StatusForbidden)
		return forbidden(c, "Specified user is suspended")
	}

	err = internal.CreateFollowRelation(currentUsr.Uuid, targetUuid)

	if err != nil {
		switch err {
		case internal.ErrCantFollowYourself:
			return badRequest(c, "Cannot follow yourself")
		case internal.ErrAlreadyFollowing:
			return badRequest(c, "You are already following")
		case internal.ErrTargetBlocked:
			return forbidden(c, "You are blocked")
		default:
			return err
		}
	}

	c.Status(fiber.StatusNoContent)
	return nil
}
