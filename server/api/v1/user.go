package v1

import (
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/MatticNote/MatticNote/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
	"strings"
	"time"
)

type (
	apiV1UserStruct struct {
		ID          string     `json:"id"`
		Username    *string    `json:"username"`
		Host        *string    `json:"host"`
		DisplayName *string    `json:"display_name"`
		Headline    *string    `json:"headline"`
		Description *string    `json:"description"`
		CreatedAt   *time.Time `json:"created_at"`
		IsModerator bool       `json:"is_moderator"`
		IsAdmin     bool       `json:"is_admin"`
	}
)

func newApiV1UserStructFromInternal(it *types.User) *apiV1UserStruct {
	s := new(apiV1UserStruct)

	s.ID = it.ID.String()
	s.IsModerator = it.IsModerator
	s.IsAdmin = it.IsAdmin

	if it.Username.Valid {
		s.Username = &it.Username.String
	}

	if it.Host.Valid {
		s.Host = &it.Host.String
	}

	if it.DisplayName.Valid {
		s.DisplayName = &it.DisplayName.String
	}

	if it.Headline.Valid {
		s.Headline = &it.Headline.String
	}

	if it.Description.Valid {
		s.Description = &it.Description.String
	}

	if it.CreatedAt.Valid {
		s.CreatedAt = &it.CreatedAt.Time
	}

	return s
}

func userApiRoute(r fiber.Router) {
	r.Get("/by/username/:acct", userGetUsername)
	r.Get("/me", loginRequired, userGetMe)
	r.Get("/:id", userGet)
}

func userGetUsername(c *fiber.Ctx) error {
	acct := strings.SplitN(c.Params("acct"), "@", 2)
	var (
		user *types.User
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

	if !user.IsActive {
		return apiGone(c, "User is gone")
	}

	return c.JSON(newApiV1UserStructFromInternal(user))
}

func userGet(c *fiber.Ctx) error {
	id, err := ksuid.Parse(c.Params("id"))
	if err != nil {
		return apiNotFound(c, "User not found")
	}

	user, err := account.GetUser(id)
	if err != nil {
		if err == account.ErrUserNotFound {
			return apiNotFound(c, "User not found")
		} else {
			return err
		}
	}

	if !user.IsActive {
		return apiGone(c, "User is gone")
	}

	return c.JSON(newApiV1UserStructFromInternal(user))
}

func userGetMe(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*types.User)

	return c.JSON(newApiV1UserStructFromInternal(user))
}
