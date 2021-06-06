package settings

import (
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strconv"
)

type editProfileFormStruct struct {
	Name           string `form:"name"`
	Summary        string `form:"summary"`
	IsBot          string `form:"is_bot"`
	AcceptManually string `form:"accept_manually"`
}

func editProfileGet(c *fiber.Ctx) error {
	return editProfileView(c, false)
}

func editProfileView(c *fiber.Ctx, isSuccess bool, errs ...string) error {
	jwtCurrentUserKey := c.Locals(internal.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	currentUser, err := internal.GetLocalUser(claim["sub"].(string))
	if err != nil {
		return err
	}

	return c.Render("account_settings/profile",
		fiber.Map{
			"IsSuccess":    isSuccess,
			"Errors":       errs,
			"CSRFFormName": misc.CSRFFormName,
			"CSRFToken":    c.Context().UserValue(misc.CSRFContextKey).(string),
			"Profile": fiber.Map{
				"Name":           currentUser.DisplayName.String,
				"Summary":        currentUser.Summary.String,
				"IsBot":          currentUser.IsBot,
				"AcceptManually": currentUser.AcceptManually,
			},
		},
		"_layout/settings",
	)
}

func editProfilePost(c *fiber.Ctx) error {
	formData := new(editProfileFormStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	jwtCurrentUserKey := c.Locals(internal.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	isBot, err := strconv.ParseBool(formData.IsBot)
	if err != nil {
		isBot = false
	}

	acceptManually, err := strconv.ParseBool(formData.AcceptManually)
	if err != nil {
		acceptManually = false
	}

	err = internal.UpdateProfile(
		uuid.MustParse(claim["sub"].(string)),
		formData.Name,
		formData.Summary,
		isBot,
		acceptManually,
	)
	if err != nil {
		return err
	}

	return editProfileView(c, true)
}
