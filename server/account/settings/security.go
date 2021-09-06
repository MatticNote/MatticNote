package settings

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/MatticNote/MatticNote/database"
	"github.com/MatticNote/MatticNote/internal/signature"
	"github.com/MatticNote/MatticNote/internal/user"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"image/png"
	"time"
)

type (
	signinLogStruct struct {
		TriedAt   time.Time
		IsSuccess bool
		FromIp    sql.NullString
	}

	user2faSetupForm struct {
		Token string `validate:"required,len=6,containsany=0123456789"`
	}

	user2faDisableFormStruct struct {
		Password string `form:"password" validate:"required"`
	}
)

func securityPageGet(c *fiber.Ctx) error {
	return securityPageView(c, false)
}

func securityPageView(c *fiber.Ctx, invalidForm bool) error {
	jwtCurrentUserKey := c.Locals(signature.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	var currentUserUuid = uuid.MustParse(claim["sub"].(string))
	is2faEnabled, err := user.StatusUser2fa(currentUserUuid)
	if err != nil {
		return err
	}
	var viewBind = fiber.Map{
		"SignInLogs":   make([]signinLogStruct, 0),
		"InvalidForm":  invalidForm,
		"Is2faEnabled": is2faEnabled,
		"CSRFFormName": misc.CSRFFormName,
		"CSRFToken":    c.Context().UserValue(misc.CSRFContextKey).(string),
	}

	rows, err := database.DBPool.Query(
		context.Background(),
		"select tried_at, is_success, from_ip from signin where target_user = $1 order by tried_at desc limit 20;",
		currentUserUuid.String(),
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Render(
				"account_settings/security",
				viewBind,
				"_layout/settings",
			)
		} else {
			return err
		}
	}
	defer rows.Close()

	for rows.Next() {
		data := signinLogStruct{}
		err = rows.Scan(
			&data.TriedAt,
			&data.IsSuccess,
			&data.FromIp,
		)
		if err != nil {
			return err
		}
		viewBind["SignInLogs"] = append(viewBind["SignInLogs"].([]signinLogStruct), data)
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	return c.Render(
		"account_settings/security",
		viewBind,
		"_layout/settings",
	)
}

func setup2faGet(c *fiber.Ctx) error {
	return setup2faView(c, false)
}

func setup2faView(c *fiber.Ctx, isFail bool) error {
	jwtCurrentUserKey := c.Locals(signature.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	totpCode, err := user.Setup2faCode(uuid.MustParse(claim["sub"].(string)))
	if err != nil {
		if err == user.Err2faAlreadyEnabled {
			return c.Redirect("/account/settings/security")
		} else {
			return err
		}
	}

	image, err := totpCode.Image(256, 256)
	if err != nil {
		return err
	}
	var qrBytes bytes.Buffer
	err = png.Encode(&qrBytes, image)
	if err != nil {
		return err
	}

	return c.Render(
		"account_settings/security_2fa_setup",
		fiber.Map{
			"QRCode":       fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(qrBytes.Bytes())),
			"SecretCode":   totpCode.Secret(),
			"CSRFFormName": misc.CSRFFormName,
			"CSRFToken":    c.Context().UserValue(misc.CSRFContextKey).(string),
			"IsFail":       isFail,
		},
		"_layout/settings",
	)
}

func setup2faPost(c *fiber.Ctx) error {
	jwtCurrentUserKey := c.Locals(signature.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	formData := new(user2faSetupForm)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(formData); errs != nil {
		return setup2faView(c, true)
	}

	targetUuid := uuid.MustParse(claim["sub"].(string))

	err := user.Validate2faCode(targetUuid, formData.Token)
	if err != nil {
		if err == user.ErrInvalid2faToken {
			return setup2faView(c, true)
		} else {
			return err
		}
	}

	if err := user.Enable2faAuth(targetUuid); err != nil {
		return err
	}

	return c.Redirect("/account/settings/security/2fa/backup")
}

func get2faBackup(c *fiber.Ctx) error {
	jwtCurrentUserKey := c.Locals(signature.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	code, err := user.Get2faBackupCode(uuid.MustParse(claim["sub"].(string)))
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Redirect("/account/settings/security/2fa/setup", fiber.StatusTemporaryRedirect)
		} else {
			return err
		}
	}

	return c.Render(
		"account_settings/security_2fa_backup",
		fiber.Map{
			"BackupCodes":  code,
			"CSRFFormName": misc.CSRFFormName,
			"CSRFToken":    c.Context().UserValue(misc.CSRFContextKey).(string),
		},
		"_layout/settings",
	)
}

func regenerate2faBackup(c *fiber.Ctx) error {
	jwtCurrentUserKey := c.Locals(signature.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	err := user.Regenerate2faBackupCode(uuid.MustParse(claim["sub"].(string)))
	if err != nil {
		return err
	}

	return get2faBackup(c)
}

func disable2faPost(c *fiber.Ctx) error {
	jwtCurrentUserKey := c.Locals(signature.JWTContextKey).(*jwt.Token)
	if !jwtCurrentUserKey.Valid {
		return fiber.ErrForbidden
	}

	claim := jwtCurrentUserKey.Claims.(jwt.MapClaims)

	formData := new(user2faDisableFormStruct)

	if err := c.BodyParser(formData); err != nil {
		return err
	}

	if errs := misc.ValidateForm(formData); errs != nil {
		return securityPageView(c, true)
	}

	targetUuid := uuid.MustParse(claim["sub"].(string))

	err := user.ValidateUserPassword(targetUuid, formData.Password)
	if err != nil {
		if err == user.ErrInvalidPassword {
			return securityPageView(c, true)
		} else {
			return err
		}
	}

	err = user.Disable2faAuth(targetUuid)
	if err != nil {
		if err == user.ErrCantDisable2fa {
			return c.Redirect("/account/settings/security")
		} else {
			return err
		}
	}

	return c.Redirect("/account/settings/security")
}
