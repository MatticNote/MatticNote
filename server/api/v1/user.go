package v1

import (
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
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

func newApiV1UserStructFromSchema(it *schemas.User) *apiV1UserStruct {
	return nil
}

func userApiRoute(r fiber.Router) {
	r.Get("/by/username/:acct", userGetUsername)
	r.Get("/me", loginRequired, userGetMe)
	r.Get("/:id", userGet)
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

	// TODO: Deleted validation

	return c.JSON(newApiV1UserStructFromSchema(user))
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

	//if !user.IsActive {
	//	return apiGone(c, "User is gone")
	//}

	return c.JSON(newApiV1UserStructFromSchema(user))
}

func userGetMe(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)

	return c.JSON(newApiV1UserStructFromSchema(user))
}
