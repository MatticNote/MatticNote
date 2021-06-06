package account

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
)

type loginUserStruct struct {
	Login    string `validate:"required"`
	Password string `validate:"required"`
}

func loginUserGet(c *fiber.Ctx) error {
	if c.Cookies(internal.JWTAuthCookieName, "") != "" {
		return c.Redirect(c.Query("next", "/web/"), 307)
	}

	return loginUserView(c)
}

func loginUserView(c *fiber.Ctx, errors ...string) error {
	return c.Status(fiber.StatusOK).Render(
		"account/login",
		fiber.Map{
			"CSRFFormName": misc.CSRFFormName,
			"CSRFToken":    c.Context().UserValue(misc.CSRFContextKey).(string),
			"errors":       errors,
		},
		"_layout/account",
	)
}

func loginPost(c *fiber.Ctx) error {
	formData := new(loginUserStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(*formData); errs != nil {
		return loginUserView(c, errs...)
	}

	targetUuid, err := internal.ValidateLoginUser(formData.Login, formData.Password)
	if err != nil {
		switch err {
		case internal.ErrLoginFailed:
			return loginUserView(c, "Incorrect login name or password")
		case internal.ErrEmailAuthRequired:
			return loginUserView(c, "Email authentication required")
		default:
			return err
		}
	}

	jwtSignedString, err := internal.SignJWT(targetUuid)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     internal.JWTAuthCookieName,
		Value:    jwtSignedString,
		Path:     "/",
		Secure:   config.Config.Server.CookieSecure,
		HTTPOnly: false,
		SameSite: "Strict",
		MaxAge:   int(internal.JWTSignExpiredDuration),
	})

	return c.Redirect(c.Query("next", "/web/"))
}
