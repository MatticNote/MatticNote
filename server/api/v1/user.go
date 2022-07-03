package v1

import (
	"errors"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
	"strings"
	"time"
)

type (
	apiV1UserStruct struct {
		ID          string    `json:"id"`
		Username    *string   `json:"username"`
		Host        *string   `json:"host"`
		DisplayName *string   `json:"display_name"`
		Headline    *string   `json:"headline"`
		Description *string   `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		IsModerator bool      `json:"is_moderator"`
		IsAdmin     bool      `json:"is_admin"`
	}

	apiV1UserUpdateStruct struct {
		DisplayName *string `json:"display_name"`
		Headline    *string `json:"headline"`
		Description *string `json:"description"`
	}

	apiV1UserFollowPendingStruct struct {
		IsPending bool `json:"is_pending"`
	}
)

func newApiV1UserStructFromSchema(it *schemas.User) *apiV1UserStruct {
	u := new(apiV1UserStruct)

	u.ID = it.ID.String()
	u.CreatedAt = it.CreatedAt
	u.IsModerator = it.IsModerator
	u.IsAdmin = it.IsAdmin

	if it.Username.Valid {
		u.Username = &it.Username.String
	}

	if it.Host.Valid {
		u.Host = &it.Host.String
	}

	if it.DisplayName.Valid {
		u.DisplayName = &it.DisplayName.String
	}

	if it.Headline.Valid {
		u.Headline = &it.Headline.String
	}

	if it.Description.Valid {
		u.Description = &it.Description.String
	}

	return u
}

func userApiRoute(r fiber.Router) {
	r.Get("/by/username/:acct", userGetUsername)
	r.Get("/me", loginRequired, userGetMe)
	r.Put("/me", loginRequired, userUpdateMe)
	r.Get("/:id", userGet)
	r.Post("/:id/follow", loginRequired, activeAccountRequired, userFollow)
	r.Delete("/:id/follow", loginRequired, activeAccountRequired, userUnfollow)
}

func userGetUsername(c *fiber.Ctx) error {
	acct := strings.SplitN(c.Params("acct"), "@", 2)
	var (
		user *schemas.User
		err  error
	)
	if len(acct) > 1 {
		user, err = account.GetUserByUsername(acct[0], acct[1])
	} else {
		user, err = account.GetUserByUsername(acct[0])
	}
	if err != nil {
		if err == account.ErrUserNotFound {
			return apiNotFound(c, "User not found")
		} else {
			return err
		}
	}

	if user.DeletedAt.Valid && user.DeletedAt.Time.Before(time.Now()) {
		return apiGone(c, "User is gone")
	}

	return c.JSON(newApiV1UserStructFromSchema(user))
}

func userGet(c *fiber.Ctx) error {
	id, err := ksuid.Parse(c.Params("id"))
	if err != nil {
		return apiNotFound(c, "User not found")
	}

	user, err := account.GetUser(id)
	if err != nil {
		if errors.Is(err, account.ErrUserNotFound) {
			return apiNotFound(c, "User not found")
		} else {
			return err
		}
	}

	if user.DeletedAt.Valid && user.DeletedAt.Time.Before(time.Now()) {
		return apiGone(c, "User is gone")
	}

	return c.JSON(newApiV1UserStructFromSchema(user))
}

func userGetMe(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)

	return c.JSON(newApiV1UserStructFromSchema(user))
}

func userUpdateMe(c *fiber.Ctx) error {
	body := new(apiV1UserUpdateStruct)

	err := c.BodyParser(body)
	if err != nil {
		return apiBadRequest(c, "Invalid form.")
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return apiBadRequest(c, "Invalid form.")
	}

	user := c.Locals("currentUser").(*schemas.User)

	if body.DisplayName != nil {
		err := user.DisplayName.Scan(*body.DisplayName)
		if err != nil {
			return err
		}
	} else {
		err := user.DisplayName.Scan(nil)
		if err != nil {
			return err
		}
	}

	if body.Headline != nil {
		err := user.Headline.Scan(*body.Headline)
		if err != nil {
			return err
		}
	} else {
		err := user.Headline.Scan(nil)
		if err != nil {
			return err
		}
	}

	if body.Description != nil {
		err := user.Description.Scan(*body.Description)
		if err != nil {
			return err
		}
	} else {
		err := user.Description.Scan(nil)
		if err != nil {
			return err
		}
	}

	err = account.UpdateUser(user)
	if err != nil {
		return err
	}

	return c.JSON(newApiV1UserStructFromSchema(user))
}

func userFollow(c *fiber.Ctx) error {
	targetId, err := ksuid.Parse(c.Params("id"))
	if err != nil {
		return apiNotFound(c, "User not found")
	}

	user := c.Locals("currentUser").(*schemas.User)

	isActive, err := account.CreateFollowRelation(user.ID, targetId)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound):
			return apiNotFound(c, "User not found")
		case errors.Is(err, account.ErrCannotRelationYourself):
			return apiForbidden(c, "You can't follow yourself")
		case errors.Is(err, account.ErrAlreadyFollowing):
			return apiBadRequest(c, "You are already following this user")
		default:
			return err
		}
	}

	return c.JSON(apiV1UserFollowPendingStruct{IsPending: !isActive})
}

func userUnfollow(c *fiber.Ctx) error {
	targetId, err := ksuid.Parse(c.Params("id"))
	if err != nil {
		return apiNotFound(c, "User not found")
	}

	user := c.Locals("currentUser").(*schemas.User)

	err = account.DeleteFollowRelation(user.ID, targetId)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrCannotRelationYourself):
			return apiForbidden(c, "You can't unfollow yourself")
		case errors.Is(err, account.ErrNotFollowing):
			return apiBadRequest(c, "You are not following this user")
		default:
			return err
		}
	}

	return c.Status(fiber.StatusNoContent).Send([]byte(""))
}
