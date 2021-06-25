package v1

import (
	"context"
	"errors"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

var (
	errUnauthorized = errors.New("unauthorized")
	errInvalidUUID  = errors.New("invalid UUID")
	errNotFoundUser = errors.New("user not found")
	errUserSuspend  = errors.New("user is suspended")
	errUserGone     = errors.New("user is gone")
)

func getCurrentUser(c *fiber.Ctx) (*internal.LocalUserStruct, error) {
	currentUsr, ok := c.Locals(internal.LoginUserLocal).(*internal.LocalUserStruct)
	if !ok {
		return nil, errUnauthorized
	}

	return currentUsr, nil
}

func validateUser(c *fiber.Ctx) error {
	targetUuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return errInvalidUUID
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
			return errNotFoundUser
		} else {
			return err
		}
	}

	if !isActive {
		return errUserGone
	}

	if isSuspend {
		return errUserSuspend
	}

	return nil
}

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
	err := validateUser(c)
	if err != nil {
		switch err {
		case errNotFoundUser:
			return notFound(c)
		case errUserSuspend:
			return forbidden(c)
		case errInvalidUUID:
			return badRequest(c, "Invalid UUID format")
		case errUserGone:
			c.Status(fiber.StatusGone)
			return nil
		default:
			return err
		}
	}

	currentUsr, err := getCurrentUser(c)
	if err != nil {
		if err == errUnauthorized {
			return unauthorized(c)
		} else {
			return err
		}
	}

	target, err := internal.GetUser(uuid.MustParse(c.Params("uuid")))
	if err != nil {
		return err
	}

	err = internal.CreateFollowRelation(currentUsr.Uuid, target.Uuid, target.AcceptManually)

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
		"is_pending": target.AcceptManually,
	})
}

func unfollowUser(c *fiber.Ctx) error {
	err := validateUser(c)
	if err != nil {
		switch err {
		case errNotFoundUser:
			return notFound(c)
		case errUserSuspend:
			return forbidden(c)
		case errInvalidUUID:
			return badRequest(c, "Invalid UUID format")
		case errUserGone:
			c.Status(fiber.StatusGone)
			return nil
		default:
			return err
		}
	}

	currentUsr, err := getCurrentUser(c)
	if err != nil {
		if err == errUnauthorized {
			return unauthorized(c)
		} else {
			return err
		}
	}

	target, err := internal.GetUser(uuid.MustParse(c.Params("uuid")))
	if err != nil {
		return err
	}

	err = internal.DestroyFollowRelation(currentUsr.Uuid, target.Uuid)

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

func muteUser(c *fiber.Ctx) error {
	err := validateUser(c)
	if err != nil {
		switch err {
		case errNotFoundUser:
			return notFound(c)
		case errUserSuspend:
			return forbidden(c)
		case errInvalidUUID:
			return badRequest(c, "Invalid UUID format")
		case errUserGone:
			c.Status(fiber.StatusGone)
			return nil
		default:
			return err
		}
	}

	currentUsr, err := getCurrentUser(c)
	if err != nil {
		if err == errUnauthorized {
			return unauthorized(c)
		} else {
			return err
		}
	}

	target, err := internal.GetUser(uuid.MustParse(c.Params("uuid")))
	if err != nil {
		return err
	}

	err = internal.CreateMuteRelation(currentUsr.Uuid, target.Uuid)

	if err != nil {
		switch err {
		case internal.ErrCantRelateYourself:
			return badRequest(c, "Cannot follow yourself")
		case internal.ErrAlreadyMuting:
			return badRequest(c, "You are already muting")
		default:
			return err
		}
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func unmuteUser(c *fiber.Ctx) error {
	err := validateUser(c)
	if err != nil {
		switch err {
		case errNotFoundUser:
			return notFound(c)
		case errUserSuspend:
			return forbidden(c)
		case errInvalidUUID:
			return badRequest(c, "Invalid UUID format")
		case errUserGone:
			c.Status(fiber.StatusGone)
			return nil
		default:
			return err
		}
	}

	currentUsr, err := getCurrentUser(c)
	if err != nil {
		if err == errUnauthorized {
			return unauthorized(c)
		} else {
			return err
		}
	}

	target, err := internal.GetUser(uuid.MustParse(c.Params("uuid")))
	if err != nil {
		return err
	}

	err = internal.DestroyMuteRelation(currentUsr.Uuid, target.Uuid)

	if err != nil {
		switch err {
		case internal.ErrCantRelateYourself:
			return badRequest(c, "Cannot follow yourself")
		case internal.ErrNotMuting:
			return badRequest(c, "You are not muting")
		default:
			return err
		}
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func blockUser(c *fiber.Ctx) error {
	err := validateUser(c)
	if err != nil {
		switch err {
		case errNotFoundUser:
			return notFound(c)
		case errUserSuspend:
			return forbidden(c)
		case errInvalidUUID:
			return badRequest(c, "Invalid UUID format")
		case errUserGone:
			c.Status(fiber.StatusGone)
			return nil
		default:
			return err
		}
	}

	currentUsr, err := getCurrentUser(c)
	if err != nil {
		if err == errUnauthorized {
			return unauthorized(c)
		} else {
			return err
		}
	}

	target, err := internal.GetUser(uuid.MustParse(c.Params("uuid")))
	if err != nil {
		return err
	}

	err = internal.CreateBlockRelation(currentUsr.Uuid, target.Uuid)

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
	err := validateUser(c)
	if err != nil {
		switch err {
		case errNotFoundUser:
			return notFound(c)
		case errUserSuspend:
			return forbidden(c)
		case errInvalidUUID:
			return badRequest(c, "Invalid UUID format")
		case errUserGone:
			c.Status(fiber.StatusGone)
			return nil
		default:
			return err
		}
	}

	currentUsr, err := getCurrentUser(c)
	if err != nil {
		if err == errUnauthorized {
			return unauthorized(c)
		} else {
			return err
		}
	}

	target, err := internal.GetUser(uuid.MustParse(c.Params("uuid")))
	if err != nil {
		return err
	}

	err = internal.DestroyBlockRelation(currentUsr.Uuid, target.Uuid)

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

func isUsedUsername(c *fiber.Ctx) error {
	username := c.Query("username", "")
	if username == "" {
		return badRequest(c, "query username cannot be empty")
	}

	isUsed, err := internal.CheckUsernameUsed(username)
	if err != nil {
		return err
	}

	return c.JSON(isUsed)
}

func relateUser(c *fiber.Ctx) error {
	err := validateUser(c)
	if err != nil {
		switch err {
		case errNotFoundUser:
			return notFound(c)
		case errUserSuspend:
			return forbidden(c)
		case errInvalidUUID:
			return badRequest(c, "Invalid UUID format")
		case errUserGone:
			c.Status(fiber.StatusGone)
			return nil
		default:
			return err
		}
	}

	currentUsr, err := getCurrentUser(c)
	if err != nil {
		if err == errUnauthorized {
			return unauthorized(c)
		} else {
			return err
		}
	}

	target, err := internal.GetUser(uuid.MustParse(c.Params("uuid")))
	if err != nil {
		return err
	}

	relation, err := internal.LookupUserRelation(currentUsr.Uuid, target.Uuid)
	if err != nil {
		if err == internal.ErrCantRelateYourself {
			return badRequest(c, "Can't lookup yourself")
		} else {
			return err
		}
	}

	return c.JSON(convFromInternalUserRelate(*relation))
}
