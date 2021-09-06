package account

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/MatticNote/MatticNote/internal/user"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"regexp"
)

type loginUserStruct struct {
	Login    string `validate:"required"`
	Password string `validate:"required"`
}

type user2faStruct struct {
	Code string `validate:"required"`
}

func loginUserGet(c *fiber.Ctx) error {
	if c.Cookies(signature.JWTAuthCookieName, "") != "" {
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

	nextQuery := c.Query("next", "/web/")
	if !regexp.MustCompile(`^/[a-zA-Z0-9\-_@].*$`).Match([]byte(nextQuery)) {
		c.Status(fiber.StatusForbidden)
		return nil
	}

	targetUuid, err := user.ValidateLoginUser(formData.Login, formData.Password)
	var isSuccess = false
	var user2faRequired = false
	defer func() {
		if targetUuid != uuid.Nil && !user2faRequired {
			_ = user.InsertSigninLog(targetUuid, c.IP(), isSuccess)
		}
	}()
	if err != nil {
		switch err {
		case user.ErrLoginFailed:
			return loginUserView(c, "Incorrect login name or password")
		case user.ErrEmailAuthRequired:
			return loginUserView(c, "Email authentication required")
		case user.Err2faRequired:
			user2faRequired = true
			s, err := login2faSession.Get(c)
			if err != nil {
				return err
			}
			s.Set("targetUuid", targetUuid.String())
			s.Set("next", c.Query("next", "/web/"))
			if err := s.Save(); err != nil {
				return err
			}
			return c.Redirect("/account/login/2fa", fiber.StatusFound)
		default:
			return err
		}
	}

	jwtSignedString, err := signature.SignJWT(targetUuid)
	if err != nil {
		return err
	}

	isSuccess = true
	c.Cookie(&fiber.Cookie{
		Name:     signature.JWTAuthCookieName,
		Value:    jwtSignedString,
		Path:     "/",
		Secure:   config.Config.Server.CookieSecure,
		HTTPOnly: false,
		SameSite: "Strict",
		MaxAge:   int(signature.JWTSignExpiredDuration),
	})

	return c.Redirect(c.Query("next", "/web/"))
}

func login2faGet(c *fiber.Ctx) error {
	return login2faView(c, false)
}

func login2faView(c *fiber.Ctx, isFail bool) error {
	s, err := login2faSession.Get(c)
	if err != nil {
		return err
	}
	_, ok := s.Get("targetUuid").(string)
	if !ok {
		return c.Redirect("/account/login", fiber.StatusFound)
	}

	return c.Render(
		"account/2fa",
		fiber.Map{
			"isFail":       isFail,
			"CSRFFormName": misc.CSRFFormName,
			"CSRFToken":    c.Context().UserValue(misc.CSRFContextKey).(string),
		},
		"_layout/account",
	)
}

func login2faPost(c *fiber.Ctx) error {
	s, err := login2faSession.Get(c)
	if err != nil {
		return err
	}
	targetUuidStr, ok := s.Get("targetUuid").(string)
	if !ok {
		return c.Redirect("/account/login", fiber.StatusFound)
	}
	targetUuid := uuid.MustParse(targetUuidStr)

	formData := new(user2faStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(*formData); errs != nil {
		return login2faView(c, true)
	}

	var isSuccess = false
	defer func() {
		_ = user.InsertSigninLog(targetUuid, c.IP(), isSuccess)
	}()

	err = user.Validate2faCode(targetUuid, formData.Code)
	if err != nil {
		if err == user.ErrInvalid2faToken {
			err = user.Use2faBackupCode(targetUuid, formData.Code)
			if err != nil {
				if err == user.ErrInvalid2faToken {
					return login2faView(c, true)
				} else {
					return err
				}
			}
		} else {
			return err
		}
	}

	next, ok := s.Get("next").(string)
	if !ok {
		next = "/web/"
	}

	jwtSignedString, err := signature.SignJWT(targetUuid)
	if err != nil {
		return err
	}

	isSuccess = true
	c.Cookie(&fiber.Cookie{
		Name:     signature.JWTAuthCookieName,
		Value:    jwtSignedString,
		Path:     "/",
		Secure:   config.Config.Server.CookieSecure,
		HTTPOnly: false,
		SameSite: "Strict",
		MaxAge:   int(signature.JWTSignExpiredDuration),
	})

	if err := s.Destroy(); err != nil {
		return err
	}

	return c.Redirect(next)
}
