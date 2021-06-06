package settings

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type editPasswordFormStruct struct {
	CurrentPassword string `form:"current_password" validate:"required"`
	NewPassword     string `form:"new_password" validate:"required,min=8"`
}

func editPasswordGet(c *fiber.Ctx) error {
	return editPasswordView(c, false)
}

func editPasswordView(c *fiber.Ctx, isSuccess bool, errs ...string) error {
	return c.Render("account_settings/password",
		fiber.Map{
			"IsSuccess":    isSuccess,
			"Errors":       errs,
			"CSRFFormName": misc.CSRFFormName,
			"CSRFToken":    c.Context().UserValue(misc.CSRFContextKey).(string),
		},
		"_layout/settings",
	)
}

func editPasswordPost(c *fiber.Ctx) error {
	jwtCurrentUserKey := c.Locals(internal.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	formData := new(editPasswordFormStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(*formData); errs != nil {
		return editPasswordView(c, false, errs...)
	}

	targetUuid := uuid.MustParse(claim["sub"].(string))

	err := internal.ValidateUserPassword(targetUuid, formData.CurrentPassword)
	if err != nil {
		if err == internal.ErrInvalidPassword {
			return editPasswordView(c, false, "current password does not match")
		} else {
			return err
		}
	}

	err = internal.ChangeUserPassword(targetUuid, formData.NewPassword)
	if err != nil {
		return err
	}

	return editPasswordView(c, true)
}
