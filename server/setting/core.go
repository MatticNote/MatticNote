package setting

import (
	"errors"
	"github.com/MatticNote/MatticNote/database/schemas"
	"github.com/MatticNote/MatticNote/internal/account"
	"github.com/MatticNote/MatticNote/server/common"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type changeEmailFormStruct struct {
	Email string `form:"email" validate:"required,email"`
}

func settingCoreGet(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)

	isEmailVerified, err := account.IsEmailVerified(user.ID)
	if err != nil {
		return err
	}
	email, err := account.GetUserEmail(user.ID)
	if err != nil {
		return err
	}
	changeEmail, ok := c.Locals("changeEmail").(string)
	if !ok {
		changeEmail = ""
	}
	invalidChangeEmailForm, ok := c.Locals("invalidChangeEmailForm").(bool)
	if !ok {
		invalidChangeEmailForm = false
	}

	return c.Render(
		"setting/core",
		fiber.Map{
			"user":                   user,
			"isUserEmailVerified":    isEmailVerified,
			"userEmail":              email,
			"csrfName":               common.CSRFFormName,
			"csrfToken":              c.Locals(common.CSRFContextKey),
			"changeEmail":            changeEmail,
			"invalidChangeEmailForm": invalidChangeEmailForm,
		},
		"setting/_layout",
	)
}

func settingCorePost(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*schemas.User)

	form := new(actionFormStruct)
	if err := c.BodyParser(form); err != nil {
		return fiber.ErrUnprocessableEntity
	}

	switch form.Action {
	case "change_email":
		formChangeEmail := new(changeEmailFormStruct)
		if err := c.BodyParser(formChangeEmail); err != nil {
			return fiber.ErrUnprocessableEntity
		}

		err := validator.New().Struct(*formChangeEmail)
		if err != nil {
			c.Locals("invalidChangeEmailForm", true)
			return settingCoreGet(c)
		}

		err = account.IssueEmailVerifyToken(user.ID, formChangeEmail.Email)
		if err != nil {
			if errors.Is(err, account.ErrEmailExists) {
				c.Locals("invalidChangeEmailForm", true)
				return settingCoreGet(c)
			} else {
				return err
			}
		}

		c.Locals("changeEmail", formChangeEmail.Email)
		return settingCoreGet(c)
	default:
		return fiber.ErrBadRequest
	}
}
