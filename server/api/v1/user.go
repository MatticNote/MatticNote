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
		isActive       bool
		isSuspend      bool
		acceptManually bool
	)

	err = database.DBPool.QueryRow(
		context.Background(),
		"select is_active, is_suspend, accept_manually from \"user\" where uuid = $1;",
		targetUuid.String(),
	).Scan(
		&isActive,
		&isSuspend,
		&acceptManually,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return notFound(c)
		} else {
			return err
		}
	}

	if !isActive {
		c.Status(fiber.StatusGone)
		return nil
	}

	if isSuspend {
		c.Status(fiber.StatusForbidden)
		return forbidden(c, "Specified user is suspended")
	}

	err = internal.CreateFollowRelation(currentUsr.Uuid, targetUuid, acceptManually)

	if err != nil {
		switch err {
		case internal.ErrCantRelateYourself:
			return badRequest(c, "Cannot follow yourself")
		case internal.ErrAlreadyFollowing:
			return badRequest(c, "You are already following")
		case internal.ErrTargetBlocked:
			return forbidden(c, "A designated user is a block relationship")
		default:
			return err
		}
	}

	return c.JSON(fiber.Map{
		"is_pending": acceptManually,
	})
}

func unfollowUser(c *fiber.Ctx) error {
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

	if err != nil {
		if err == pgx.ErrNoRows {
			return notFound(c)
		} else {
			return err
		}
	}

	if !isActive {
		c.Status(fiber.StatusGone)
		return nil
	}

	if isSuspend {
		c.Status(fiber.StatusForbidden)
		return forbidden(c, "Specified user is suspended")
	}

	err = internal.DestroyFollowRelation(currentUsr.Uuid, targetUuid)

	if err != nil {
		switch err {
		case internal.ErrCantRelateYourself:
			return badRequest(c, "Cannot follow yourself")
		case internal.ErrNotFollowing:
			return badRequest(c, "You are not following")
		default:
			return err
		}
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func blockUser(c *fiber.Ctx) error {
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

	if err != nil {
		if err == pgx.ErrNoRows {
			return notFound(c)
		} else {
			return err
		}
	}

	if !isActive {
		c.Status(fiber.StatusGone)
		return nil
	}

	if isSuspend {
		c.Status(fiber.StatusForbidden)
		return forbidden(c, "Specified user is suspended")
	}

	err = internal.CreateBlockRelation(currentUsr.Uuid, targetUuid)

	if err != nil {
		switch err {
		case internal.ErrCantRelateYourself:
			return badRequest(c, "Cannot follow yourself")
		case internal.ErrAlreadyBlocking:
			return badRequest(c, "You are already blocking")
		default:
			return err
		}
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func unblockUser(c *fiber.Ctx) error {
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

	if err != nil {
		if err == pgx.ErrNoRows {
			return notFound(c)
		} else {
			return err
		}
	}

	if !isActive {
		c.Status(fiber.StatusGone)
		return nil
	}

	if isSuspend {
		c.Status(fiber.StatusForbidden)
		return forbidden(c, "Specified user is suspended")
	}

	err = internal.DestroyBlockRelation(currentUsr.Uuid, targetUuid)

	if err != nil {
		switch err {
		case internal.ErrCantRelateYourself:
			return badRequest(c, "Cannot follow yourself")
		case internal.ErrNotBlocking:
			return badRequest(c, "You are not blocking")
		default:
			return err
		}
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func listRequests(c *fiber.Ctx) error {
	currentUsr, ok := c.Locals(internal.LoginUserLocal).(*internal.LocalUserStruct)
	if !ok {
		return unauthorized(c)
	}

	requests, err := internal.ListFollowRequests(currentUsr.Uuid)
	if err != nil {
		return err
	}

	// TODO: 今はUUIDだけだが、ユーザー情報とか入れたほうが良さそう
	return c.JSON(requests)
}

func acceptRequests(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return badRequest(c, "Not valid UUID format")
	}

	currentUsr, ok := c.Locals(internal.LoginUserLocal).(*internal.LocalUserStruct)
	if !ok {
		return unauthorized(c)
	}

	err = internal.AcceptFollowRequest(currentUsr.Uuid, targetUuid)
	if err != nil {
		switch err {
		case internal.ErrCantRelateYourself:
			return badRequest(c, "Can't accept yourself")
		case internal.ErrUnknownRequest:
			return badRequest(c, "Unknown follow request")
		default:
			return err
		}
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func rejectRequests(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return badRequest(c, "Not valid UUID format")
	}

	currentUsr, ok := c.Locals(internal.LoginUserLocal).(*internal.LocalUserStruct)
	if !ok {
		return unauthorized(c)
	}

	err = internal.RejectFollowRequest(currentUsr.Uuid, targetUuid)
	if err != nil {
		switch err {
		case internal.ErrCantRelateYourself:
			return badRequest(c, "Can't accept yourself")
		case internal.ErrUnknownRequest:
			return badRequest(c, "Unknown follow request")
		default:
			return err
		}
	}

	c.Status(fiber.StatusNoContent)
	return nil
}
