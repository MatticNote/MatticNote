package setting

import (
	"errors"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/MatticNote/MatticNote/server/common"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
)

type (
	changePasswordFormStruct struct {
		CurrentPassword string `form:"current_password" validate:"required"`
		NewPassword     string `form:"new_password" validate:"required,min=8"`
	}

	deleteTokenFormStruct struct {
		ID string `form:"id" validate:"required"`
	}
)

func settingSecurityGet(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)
	token := c.Cookies(common.TokenCookieName)

	tokenList, err := account.ListUserToken(user.ID)
	if err != nil {
		return err
	}

	invalidPassword, ok := c.Locals("invalidPassword").(bool)
	if !ok {
		invalidPassword = false
	}

	changedPassword, ok := c.Locals("changedPassword").(bool)
	if !ok {
		changedPassword = false
	}

	return c.Render(
		"setting/security",
		fiber.Map{
			"user":            user,
			"tokenList":       tokenList,
			"currentToken":    token,
			"invalidPassword": invalidPassword,
			"changedPassword": changedPassword,
			"csrfName":        common.CSRFFormName,
			"csrfToken":       c.Locals(common.CSRFContextKey),
		},
		"setting/_layout",
	)
}

func settingSecurityPost(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)

	form := new(actionFormStruct)
	if err := c.BodyParser(form); err != nil {
		return fiber.ErrUnprocessableEntity
	}

	switch form.Action {
	case "change_password":
		formChangePassword := new(changePasswordFormStruct)
		if err := c.BodyParser(formChangePassword); err != nil {
			return fiber.ErrUnprocessableEntity
		}

		if err := validator.New().Struct(*formChangePassword); err != nil {
			c.Locals("invalidPassword", true)
			return settingSecurityGet(c)
		}

		email, err := account.GetUserEmail(user.ID)
		if err != nil {
			return err
		}

		_, err = account.AuthenticateUser(email, formChangePassword.CurrentPassword)
		if err != nil {
			if errors.Is(err, account.ErrInvalidCredentials) {
				c.Locals("invalidPassword", true)
				return settingSecurityGet(c)
			} else {
				return err
			}
		}

		if err := account.UpdateUserPassword(user.ID, formChangePassword.NewPassword); err != nil {
			return err
		}

		c.Locals("changedPassword", true)
		return settingSecurityGet(c)
	case "revoke_session":
		formDelete := new(deleteTokenFormStruct)
		if err := c.BodyParser(formDelete); err != nil {
			return fiber.ErrUnprocessableEntity
		}

		if err := validator.New().Struct(*formDelete); err != nil {
			return fiber.ErrBadRequest
		}

		id, err := ksuid.Parse(formDelete.ID)
		if err != nil {
			return fiber.ErrBadRequest
		}

		token, err := account.GetUserToken(id)
		if err != nil {
			if errors.Is(err, account.ErrUserTokenNotFound) {
				return fiber.ErrBadRequest
			} else {
				return err
			}
		}

		if token.UserId == nil || (token.UserId != nil && user.ID != *token.UserId) {
			return fiber.ErrForbidden
		}

		err = account.DeleteUserToken(id)
		if err != nil {
			return err
		}

		return settingSecurityGet(c)
	default:
		return fiber.ErrBadRequest
	}
}
